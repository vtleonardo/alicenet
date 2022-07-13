package executor

import (
	"context"
	"sync"

	"github.com/alicenet/alicenet/consensus/db"
	"github.com/alicenet/alicenet/constants"
	"github.com/alicenet/alicenet/layer1"
	"github.com/alicenet/alicenet/layer1/executor/marshaller"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dkg"
	"github.com/alicenet/alicenet/layer1/executor/tasks/snapshots"
	monitorInterfaces "github.com/alicenet/alicenet/layer1/monitor/interfaces"
	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/alicenet/alicenet/logging"
)

type TaskSender interface {
	SendRequest(ctx context.Context, request tasks.Request) error
	Start()
	Close()
}

type SchedulerResponse struct {
	// adds a new struct where you can listen for the task response
	Err error
}

// A response channel is basically a non-blocking channel that can only be
// written and closed once.
type SchedulerResponseChannel struct {
	writeOnce sync.Once
	channel   chan *SchedulerResponse // internal channel
}

// Create a new response channel.
func NewResponseChannel() *SchedulerResponseChannel {
	return &SchedulerResponseChannel{channel: make(chan *SchedulerResponse, 1)}
}

// send a unique response and close the internal channel. Additional calls to
// this function will be no-op
func (rc *SchedulerResponseChannel) sendResponse(response *SchedulerResponse) {
	rc.writeOnce.Do(func() {
		rc.channel <- response
		close(rc.channel)
	})
}

type internalRequest struct {
	request  tasks.Request
	response *SchedulerResponseChannel
}

var _ TaskSender = &Sender{}

type Sender struct {
	taskScheduler  *TasksSchedulerBackend
	requestChannel chan internalRequest
}

func NewTasksScheduler(database *db.Database, eth layer1.Client, adminHandler monitorInterfaces.AdminHandler, txWatcher *transaction.FrontWatcher) (TaskSender, error) {

	// main context that will cancel all workers and go routine
	mainCtx, cf := context.WithCancel(context.Background())

	// Setup tasks scheduler
	taskRequestChan := make(chan internalRequest, constants.TaskSchedulerBufferSize)
	defer close(taskRequestChan)

	s := &TasksSchedulerBackend{
		Schedule:         make(map[string]TaskRequestInfo),
		mainCtx:          mainCtx,
		mainCtxCf:        cf,
		database:         database,
		eth:              eth,
		adminHandler:     adminHandler,
		marshaller:       GetTaskRegistry(),
		cancelChan:       make(chan bool, 1),
		taskRequestChan:  taskRequestChan,
		taskResponseChan: &taskResponseChan{trChan: make(chan tasks.Response, 100)},
		txWatcher:        txWatcher,
	}

	logger := logging.GetLogger("tasks")
	s.logger = logger.WithField("Component", "schedule")

	tasksManager, err := NewTaskProcessor(txWatcher, database, logger.WithField("Component", "manager"))
	if err != nil {
		return nil, err
	}
	s.tasksManager = tasksManager

	sender := &Sender{
		taskScheduler:  s,
		requestChannel: taskRequestChan,
	}

	return sender, nil
}

func (r *Sender) Start() {
	r.taskScheduler.Start()
}

func (r *Sender) Close() {
	r.taskScheduler.Close()
	close(r.requestChannel)
}

func (r *Sender) SendRequest(ctx context.Context, request tasks.Request) error {
	req := internalRequest{request: request, response: NewResponseChannel()}

	// wait for request to be accepted
	select {
	case r.requestChannel <- req:
	case <-ctx.Done():
		return ctx.Err()
	}

	// wait for request to be processed
	select {
	case response := <-req.response.channel:
		return response.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func GetTaskRegistry() *marshaller.TypeRegistry {
	// registry the type here
	tr := &marshaller.TypeRegistry{}
	tr.RegisterInstanceType(&dkg.CompletionTask{})
	tr.RegisterInstanceType(&dkg.DisputeShareDistributionTask{})
	tr.RegisterInstanceType(&dkg.DisputeMissingShareDistributionTask{})
	tr.RegisterInstanceType(&dkg.DisputeMissingKeySharesTask{})
	tr.RegisterInstanceType(&dkg.DisputeMissingGPKjTask{})
	tr.RegisterInstanceType(&dkg.DisputeGPKjTask{})
	tr.RegisterInstanceType(&dkg.GPKjSubmissionTask{})
	tr.RegisterInstanceType(&dkg.KeyShareSubmissionTask{})
	tr.RegisterInstanceType(&dkg.MPKSubmissionTask{})
	tr.RegisterInstanceType(&dkg.RegisterTask{})
	tr.RegisterInstanceType(&dkg.DisputeMissingRegistrationTask{})
	tr.RegisterInstanceType(&dkg.ShareDistributionTask{})
	tr.RegisterInstanceType(&snapshots.SnapshotTask{})
	return tr
}
