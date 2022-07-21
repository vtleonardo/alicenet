package executor

import (
	"context"
	"errors"
	"github.com/alicenet/alicenet/layer1/executor/marshaller"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"sync"
)

//////////////////////////////////////////////////////////////////////////////////////////////
//      Shared models: used for the communication between the Handler and its clients       //
//////////////////////////////////////////////////////////////////////////////////////////////

// TaskAction is an enumeration indicating the actions that the scheduler
// can do with a task during a request:
type TaskAction int

// The possible actions that the scheduler can do with a task during a request:
// * KillByType          - To kill/prune a task by type immediately
// * KillById            - To kill/prune a task by id immediately
// * Schedule            - To schedule a new task
const (
	KillByType TaskAction = iota
	KillById
	Schedule
)

func (action TaskAction) String() string {
	return [...]string{
		"KillByType",
		"KillById",
		"Schedule",
	}[action]
}

// HandlerResponse returned from the Handler to the external clients
type HandlerResponse struct {
	doneChan chan struct{}
	err      error // error in case the task failed
}

// newHandlerResponse creates HandlerResponse
func newHandlerResponse() *HandlerResponse {
	return &HandlerResponse{doneChan: make(chan struct{})}
}

// IsReady function to check if the response from the task is ready to share with Handler client
func (r *HandlerResponse) IsReady() bool {
	select {
	case <-r.doneChan:
		return true
	default:
		return false
	}
}

// GetResponseBlocking blocking function to get the execution status of a task.
// This function will block until the task is finished and the final result is returned.
func (r *HandlerResponse) GetResponseBlocking(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.doneChan:
		return r.err
	}
}

// writeResponse function to write the response or error from the task request being executed
func (r *HandlerResponse) writeResponse(err error) {
	if !r.IsReady() {
		r.err = err
		close(r.doneChan)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
//     Internal models: used for the communication between the Handler and the Manager      //
//////////////////////////////////////////////////////////////////////////////////////////////

// InternalTaskState is an enumeration indicating the possible states of a task
type InternalTaskState int

const (
	NotStarted InternalTaskState = iota
	Running
	Killed
)

func (state InternalTaskState) String() string {
	return [...]string{
		"NotStarted",
		"Running",
		"Killed",
	}[state]
}

// managerRequest used to make a request from Handler to TaskManager for scheduling or killing a task
type managerRequest struct {
	task     tasks.Task
	id       string
	action   TaskAction
	response *ManagerResponseChannel
}

// BaseRequest information needed to start managing a task request
type BaseRequest struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	Start         uint64            `json:"start"`
	End           uint64            `json:"end"`
	InternalState InternalTaskState `json:"internalState"`
}

// ManagerRequestInfo used for controlling, recovering and managing a task request.
// It's a unified struct that contains the request TaskManager receives from the Handler,
// the data is used in order to manage the task and the response from the TaskExecutor.
type ManagerRequestInfo struct {
	BaseRequest
	Task             tasks.Task
	HandlerResponse  *HandlerResponse
	ExecutorResponse ExecutorResponse
	killedAt         uint64
}

// requestStored with an internal wrapper for the task interface
type requestStored struct {
	BaseRequest
	ExecutorResponse
	WrappedTask *marshaller.InstanceWrapper `json:"wrappedTask"`
}

// innerSequentialSchedule used to store requestStored for recovery
type innerSequentialSchedule struct {
	Schedule map[string]*requestStored
}

// managerResponse is used to communicate the task response from TaskManager to Handler
type managerResponse struct {
	HandlerResponse *HandlerResponse
	Err             error
}

// ManagerResponseChannel is a non-blocking channel that can only be written and closed once.
// It's used for communication between TaskManager and Handler
type ManagerResponseChannel struct {
	writeOnce sync.Once
	channel   chan *managerResponse // internal channel
}

// NewManagerResponseChannel creates ManagerResponseChannel.
func NewManagerResponseChannel() *ManagerResponseChannel {
	return &ManagerResponseChannel{channel: make(chan *managerResponse, 1)}
}

// sendResponse and closes the internal ManagerResponseChannel.
//Additional calls to this function will be no-op
func (rc *ManagerResponseChannel) sendResponse(response *managerResponse) {
	rc.writeOnce.Do(func() {
		rc.channel <- response
		close(rc.channel)
	})
}

// listen until the response is received
func (rc *ManagerResponseChannel) listen(ctx context.Context) (*HandlerResponse, error) {
	// wait for request to be processed
	select {
	case response := <-rc.channel:
		return response.HandlerResponse, response.Err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
//     Internal models: used for the communication between the Manager and the Executor     //
//////////////////////////////////////////////////////////////////////////////////////////////

// ExecutorResponse used inside executorResponseChan to store the task execution result
type ExecutorResponse struct {
	Id  string `json:"-"`
	Err error  `json:"response_error"`
}

// executorResponseChan is used to communicate the task execution result from TaskExecutor to
// TaskManager. It can only be written and closed once.
type executorResponseChan struct {
	writeOnce sync.Once
	erChan    chan ExecutorResponse
	isClosed  bool
}

// close executorResponseChan internal erChan
func (tr *executorResponseChan) close() {
	tr.writeOnce.Do(func() {
		tr.isClosed = true
		close(tr.erChan)
	})
}

// Add ExecutorResponse to internal erChan
func (tr *executorResponseChan) Add(id string, err error) {
	if !tr.isClosed {
		tr.erChan <- ExecutorResponse{Id: id, Err: err}
	}
}

var _ tasks.InternalTaskResponseChan = &executorResponseChan{}

var (
	ErrNotScheduled                   = errors.New("scheduled task not found")
	ErrWrongParams                    = errors.New("wrong start/end height for the task")
	ErrTaskExpired                    = errors.New("the task is already expired")
	ErrTaskNotAllowMultipleExecutions = errors.New("a task of the same type is already scheduled and allowed multiple execution for this type is false")
	ErrTaskIsNil                      = errors.New("the task in the request is nil")
	ErrTaskTypeNotInRegistry          = errors.New("the task type is not in registry")
	ErrTaskIdEmpty                    = errors.New("the task id is empty")
)