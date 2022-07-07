package tasks

import (
	"context"
	"errors"

	"github.com/alicenet/alicenet/consensus/db"
	"github.com/alicenet/alicenet/layer1"
	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/sirupsen/logrus"
)

type BaseTask struct {
	// Task name/type
	Name string `json:"name"`
	// If this task can be executed in parallel with other tasks of the same type/name
	AllowMultiExecution bool `json:"allowMultiExecution"`
	// Subscription options (if the task should be retried, finality delay, etc)
	SubscribeOptions *transaction.SubscribeOptions `json:"subscribeOptions,omitempty"`
	// Unique Id of the task
	Id string `json:"id"`
	// Which block the task should be started. In case the start is 0 the task is started immediately.
	Start uint64 `json:"start"`
	// Which block the task should be ended. In case the end is 0 the task runs
	// forever (until the task succeeds or its killed, be careful when using this).
	// Otherwise, the task will end at the specified block.
	End uint64 `json:"end"`

	isInitialized    bool               `json:"-"`
	wasKilled        bool               `json:"-"`
	ctx              context.Context    `json:"-"`
	cancelFunc       context.CancelFunc `json:"-"`
	database         *db.Database       `json:"-"`
	logger           *logrus.Entry      `json:"-"`
	client           layer1.Client      `json:"-"`
	taskResponseChan TaskResponseChan   `json:"-"`
}

func NewBaseTask(start uint64, end uint64, allowMultiExecution bool, subscribeOptions *transaction.SubscribeOptions) *BaseTask {
	ctx, cf := context.WithCancel(context.Background())

	return &BaseTask{
		Start:               start,
		End:                 end,
		AllowMultiExecution: allowMultiExecution,
		SubscribeOptions:    subscribeOptions,
		ctx:                 ctx,
		cancelFunc:          cf,
	}
}

// Initialize default implementation for the ITask interface
func (bt *BaseTask) Initialize(ctx context.Context, cancelFunc context.CancelFunc, database *db.Database, logger *logrus.Entry, eth layer1.Client, name string, id string, taskResponseChan TaskResponseChan) error {
	if !bt.isInitialized {
		bt.Name = name
		bt.Id = id
		bt.ctx = ctx
		bt.cancelFunc = cancelFunc
		bt.database = database
		bt.logger = logger
		bt.client = eth
		bt.taskResponseChan = taskResponseChan
		bt.isInitialized = true
		return nil
	}
	return errors.New("trying to initialize task twice!")
}

// GetId default implementation for the ITask interface
func (bt *BaseTask) GetId() string {
	return bt.Id
}

// GetStart default implementation for the ITask interface
func (bt *BaseTask) GetStart() uint64 {
	return bt.Start
}

// GetEnd default implementation for the ITask interface
func (bt *BaseTask) GetEnd() uint64 {
	return bt.End
}

// GetName default implementation for the ITask interface
func (bt *BaseTask) GetName() string {
	return bt.Name
}

// GetAllowMultiExecution default implementation for the ITask interface
func (bt *BaseTask) GetAllowMultiExecution() bool {
	return bt.AllowMultiExecution
}

func (bt *BaseTask) GetSubscribeOptions() *transaction.SubscribeOptions {
	return bt.SubscribeOptions
}

// GetCtx default implementation for the ITask interface
func (bt *BaseTask) GetCtx() context.Context {
	return bt.ctx
}

// GetCtx default implementation for the ITask interface
func (bt *BaseTask) WasKilled() bool {
	return bt.wasKilled
}

// GetEth default implementation for the ITask interface
func (bt *BaseTask) GetClient() layer1.Client {
	return bt.client
}

// GetLogger default implementation for the ITask interface
func (bt *BaseTask) GetLogger() *logrus.Entry {
	return bt.logger
}

// Close default implementation for the ITask interface
func (bt *BaseTask) Close() {
	if bt.cancelFunc != nil {
		bt.cancelFunc()
	}
	bt.wasKilled = true
}

// Finish default implementation for the ITask interface
func (bt *BaseTask) Finish(err error) {
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			bt.logger.WithError(err).Error("got an error when executing task")
		} else {
			bt.logger.WithError(err).Debug("cancelling task execution")
		}
	} else {
		bt.logger.Info("task is done")
	}
	if bt.taskResponseChan != nil {
		bt.taskResponseChan.Add(Response{Id: bt.Id, Err: err})
	}
}

func (bt *BaseTask) GetDB() *db.Database {
	return bt.database
}
