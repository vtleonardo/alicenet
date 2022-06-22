package executor

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/layer1"
	"github.com/MadBase/MadNet/layer1/executor/marshaller"
	"github.com/MadBase/MadNet/layer1/executor/tasks"
	"github.com/MadBase/MadNet/layer1/executor/tasks/dkg"
	"github.com/MadBase/MadNet/layer1/executor/tasks/snapshots"
	monitorInterfaces "github.com/MadBase/MadNet/layer1/monitor/interfaces"
	"github.com/MadBase/MadNet/layer1/transaction"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrNotScheduled = errors.New("scheduled task not found")
	ErrWrongParams  = errors.New("wrong start/end height for the task")
	ErrTaskExpired  = errors.New("the task is already expired")
)

const (
	heightToleranceBeforeRemoving uint64 = 50
)

type TaskRequestInfo struct {
	Id        string     `json:"id"`
	Start     uint64     `json:"start"`
	End       uint64     `json:"end"`
	Task      tasks.Task `json:"-"`
	isRunning bool       `json:"-"`
}

type taskRequestInner struct {
	Id          string
	Start       uint64
	End         uint64
	WrappedTask *marshaller.InstanceWrapper
}

type TasksScheduler struct {
	Schedule         map[string]TaskRequestInfo     `json:"schedule"`
	LastHeightSeen   uint64                         `json:"last_height_seen"`
	eth              layer1.Client                  `json:"-"`
	database         *db.Database                   `json:"-"`
	adminHandler     monitorInterfaces.AdminHandler `json:"-"`
	marshaller       *marshaller.TypeRegistry       `json:"-"`
	cancelChan       chan bool                      `json:"-"`
	taskRequestChan  <-chan tasks.Task              `json:"-"`
	taskResponseChan *taskResponseChan              `json:"-"`
	taskKillChan     <-chan string                  `json:"-"`
	logger           *logrus.Entry                  `json:"-"`
	tasksManager     *TasksManager                  `json:"-"`
	txWatcher        *transaction.FrontWatcher      `json:"-"`
}

type taskResponseChan struct {
	writeOnce sync.Once
	trChan    chan tasks.TaskResponse
	isClosed  bool
}

func (tr *taskResponseChan) close() {
	tr.writeOnce.Do(func() {
		tr.isClosed = true
		close(tr.trChan)
	})
}

func (tr *taskResponseChan) Add(taskResponse tasks.TaskResponse) {
	if !tr.isClosed {
		tr.trChan <- taskResponse
	}
}

var _ tasks.TaskResponseChan = &taskResponseChan{}

type innerSequentialSchedule struct {
	Schedule map[string]*taskRequestInner
}

func GetTaskLogger(task tasks.Task) *logrus.Entry {
	logger := logging.GetLogger("tasks")
	logEntry := logger.WithFields(logrus.Fields{
		"Component": "task",
		"taskId":    task.GetId(),
		"taskName":  task.GetName(),
		"taskStart": task.GetStart(),
		"taskEnd":   task.GetEnd(),
	})
	return logEntry
}

func NewTasksScheduler(database *db.Database, eth layer1.Client, adminHandler monitorInterfaces.AdminHandler, taskRequestChan <-chan tasks.Task, taskKillChan <-chan string, txWatcher *transaction.FrontWatcher) (*TasksScheduler, error) {
	tr := &marshaller.TypeRegistry{}
	///////////////////// Add new tasks types here /////////////////////////
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
	//////////////////////////////////////////////////////////////////////////

	s := &TasksScheduler{
		Schedule:         make(map[string]TaskRequestInfo),
		database:         database,
		eth:              eth,
		adminHandler:     adminHandler,
		marshaller:       tr,
		cancelChan:       make(chan bool, 1),
		taskRequestChan:  taskRequestChan,
		taskResponseChan: &taskResponseChan{trChan: make(chan tasks.TaskResponse, 100)},
		taskKillChan:     taskKillChan,
		txWatcher:        txWatcher,
	}

	logger := logging.GetLogger("tasks")
	s.logger = logger.WithField("Component", "schedule")

	tasksManager, err := NewTaskManager(txWatcher, database, logger.WithField("Component", "manager"))
	if err != nil {
		return nil, err
	}
	s.tasksManager = tasksManager

	return s, nil
}

func (s *TasksScheduler) Start() error {
	err := s.loadState()
	if err != nil {
		s.logger.Warnf("could not find previous State: %v", err)
		if err != badger.ErrKeyNotFound {
			return err
		}
	}

	s.logger.Info(strings.Repeat("-", 80))
	s.logger.Infof("Current Tasks: %d", len(s.Schedule))
	for id, task := range s.Schedule {
		s.logger.Infof("...ID: %s Name: %s Between: %d and %d", id, task.Task.GetName(), task.Start, task.End)
	}
	s.logger.Info(strings.Repeat("-", 80))

	go s.eventLoop()
	return nil
}

func (s *TasksScheduler) Close() {
	s.logger.Warn("Closing scheduler")
	s.cancelChan <- true
}

func (s *TasksScheduler) eventLoop() {
	ctx, cf := context.WithCancel(context.Background())
	processingTime := time.After(constants.TaskSchedulerProcessingTime)

	for {
		select {
		case <-s.cancelChan:
			s.logger.Warn("Received cancel request for event loop.")
			cf()
			s.taskResponseChan.close()
			return
		case taskRequest := <-s.taskRequestChan:
			s.logger.Trace("received request for a task")
			err := s.schedule(ctx, taskRequest)
			if err != nil {
				// if we are not synchronized, don't log expired task as errors, since we will
				// be replaying the events from far way in the past
				if errors.Is(err, ErrTaskExpired) && !s.adminHandler.IsSynchronized() {
					s.logger.WithError(err).Debugf("Failed to schedule task request %d", s.LastHeightSeen)
				} else {
					s.logger.WithError(err).Errorf("Failed to schedule task request %d", s.LastHeightSeen)
				}
			}
			err = s.persistState()
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to persist state %d on task request", s.LastHeightSeen)
			}
		case taskResponse := <-s.taskResponseChan.trChan:
			s.logger.Trace("received a task response")
			err := s.processTaskResponse(ctx, taskResponse)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to processTaskResponse %v", taskResponse)
			}
			err = s.persistState()
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to persist state %d on task response", s.LastHeightSeen)
			}
		case taskToKill := <-s.taskKillChan:
			s.logger.Trace("received request to kill a task")
			err := s.killTaskByName(ctx, taskToKill)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to killTaskByName %v", taskToKill)
			}
		case <-processingTime:
			s.logger.Trace("processing latest height")
			networkCtx, networkCf := context.WithTimeout(ctx, constants.TaskSchedulerNetworkTimeout)
			height, err := s.eth.GetFinalizedHeight(networkCtx)
			networkCf()
			if err != nil {
				s.logger.WithError(err).Debug("Failed to retrieve the latest height from eth node")
				continue
			}
			s.LastHeightSeen = height

			toStart, expired, unresponsive := s.findTasks()
			err = s.startTasks(ctx, toStart)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to startTasks %d", s.LastHeightSeen)
			}
			err = s.persistState()
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to persist state %d", s.LastHeightSeen)
			}

			err = s.killTasks(ctx, expired)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to killExpiredTasks %d", s.LastHeightSeen)
			}

			err = s.removeUnresponsiveTasks(ctx, unresponsive)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to removeUnresponsiveTasks %d", s.LastHeightSeen)
			}
			err = s.persistState()
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to persist state %d", s.LastHeightSeen)
			}
			processingTime = time.After(constants.TaskSchedulerProcessingTime)
		}
	}
}

func (s *TasksScheduler) schedule(ctx context.Context, task tasks.Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		start := task.GetStart()
		end := task.GetEnd()

		if start != 0 && end != 0 && start >= end {
			return ErrWrongParams
		}

		if end != 0 && end <= s.LastHeightSeen {
			return ErrTaskExpired
		}

		id := uuid.New()
		s.Schedule[id.String()] = TaskRequestInfo{Id: id.String(), Start: start, End: end, Task: task}
		GetTaskLogger(task).Debug("Received task request")
	}
	return nil
}

func (s *TasksScheduler) processTaskResponse(ctx context.Context, taskResponse tasks.TaskResponse) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		logger := s.logger
		task, present := s.Schedule[taskResponse.Id]
		if present {
			logger = GetTaskLogger(task.Task)
		}
		if taskResponse.Err != nil {
			logger.Errorf("Task id: %s executed with error: %v", taskResponse.Id, taskResponse.Err)
		} else {
			logger.Infof("Task id: %s successfully executed", taskResponse.Id)
		}
		err := s.remove(taskResponse.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TasksScheduler) startTasks(ctx context.Context, tasks []TaskRequestInfo) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.logger.Debug("Looking for starting tasks")
		for i := 0; i < len(tasks); i++ {
			task := tasks[i]
			logEntry := GetTaskLogger(task.Task)
			logEntry.Info("task is about to start")

			go s.tasksManager.ManageTask(ctx, task.Task, task.Id, s.database, logEntry, s.eth, s.taskResponseChan)

			task.isRunning = true
			s.Schedule[task.Id] = task
		}

	}

	return nil
}

func (s *TasksScheduler) killTaskByName(ctx context.Context, taskName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.logger.Debugf("Looking for killing tasks by name %s", taskName)
		return s.killTasks(ctx, s.findTasksByName(taskName))
	}
}

func (s *TasksScheduler) killTasks(ctx context.Context, tasks []TaskRequestInfo) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for i := 0; i < len(tasks); i++ {
			task := tasks[i]
			GetTaskLogger(task.Task).Info("Task is about to be killed")
			task.Task.Close()
		}
	}

	return nil
}

func (s *TasksScheduler) removeUnresponsiveTasks(ctx context.Context, tasks []TaskRequestInfo) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.logger.Debug("Looking for removing unresponsive tasks")

		for i := 0; i < len(tasks); i++ {
			task := tasks[i]
			GetTaskLogger(task.Task).Info("Task is about to be removed for being unresponsive or expired")

			err := s.remove(task.Id)
			if err != nil {
				s.logger.WithError(err).Errorf("Failed to remove unresponsive task id: %s", task.Id)
			}
		}

	}

	return nil
}

func (s *TasksScheduler) findTasks() ([]TaskRequestInfo, []TaskRequestInfo, []TaskRequestInfo) {
	toStart := make([]TaskRequestInfo, 0)
	expired := make([]TaskRequestInfo, 0)
	unresponsive := make([]TaskRequestInfo, 0)

	for _, taskRequest := range s.Schedule {
		if taskRequest.End != 0 && taskRequest.End+heightToleranceBeforeRemoving <= s.LastHeightSeen {
			unresponsive = append(unresponsive, taskRequest)
			continue
		}

		if taskRequest.End != 0 && taskRequest.End <= s.LastHeightSeen {
			expired = append(expired, taskRequest)
			continue
		}

		if ((taskRequest.Start == 0 && taskRequest.End == 0) ||
			(taskRequest.Start != 0 && taskRequest.Start <= s.LastHeightSeen && taskRequest.End == 0) ||
			(taskRequest.Start <= s.LastHeightSeen && taskRequest.End > s.LastHeightSeen)) && !taskRequest.isRunning {

			if taskRequest.Task.GetAllowMultiExecution() ||
				(!taskRequest.Task.GetAllowMultiExecution() && len(s.findRunningTasksByName(taskRequest.Task.GetName())) == 0) {
				toStart = append(toStart, taskRequest)
			}
			continue
		}
	}
	return toStart, expired, unresponsive
}

func (s *TasksScheduler) findTasksByName(taskName string) []TaskRequestInfo {
	tasks := make([]TaskRequestInfo, 0)

	for _, taskRequest := range s.Schedule {
		if taskRequest.Task.GetName() == taskName {
			tasks = append(tasks, taskRequest)
		}
	}
	return tasks
}

func (s *TasksScheduler) findRunningTasksByName(taskName string) []TaskRequestInfo {
	tasks := make([]TaskRequestInfo, 0)

	for _, taskRequest := range s.Schedule {
		if taskRequest.Task.GetName() == taskName && taskRequest.isRunning {
			tasks = append(tasks, taskRequest)
		}
	}
	return tasks
}

func (s *TasksScheduler) length() int {
	return len(s.Schedule)
}

func (s *TasksScheduler) remove(id string) error {
	_, present := s.Schedule[id]
	if !present {
		return ErrNotScheduled
	}

	delete(s.Schedule, id)

	return nil
}

func (s *TasksScheduler) persistState() error {
	rawData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = s.database.Update(func(txn *badger.Txn) error {
		key := dbprefix.PrefixTaskSchedulerState()
		s.logger.WithField("Key", string(key)).Debug("Saving state")
		if err := utils.SetValue(txn, key, rawData); err != nil {
			s.logger.Error("Failed to set Value")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := s.database.Sync(); err != nil {
		s.logger.Error("Failed to set sync")
		return err
	}

	return nil
}

func (s *TasksScheduler) loadState() error {

	if err := s.database.View(func(txn *badger.Txn) error {
		key := dbprefix.PrefixTaskSchedulerState()
		s.logger.WithField("Key", string(key)).Debug("Looking up state")
		rawData, err := utils.GetValue(txn, key)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rawData, s)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// synchronizing db state to disk
	if err := s.database.Sync(); err != nil {
		s.logger.Error("Failed to set sync")
		return err
	}

	return nil

}

func (s *TasksScheduler) MarshalJSON() ([]byte, error) {

	ws := &innerSequentialSchedule{Schedule: make(map[string]*taskRequestInner)}

	for k, v := range s.Schedule {
		wt, err := s.marshaller.WrapInstance(v.Task)
		if err != nil {
			return []byte{}, err
		}
		ws.Schedule[k] = &taskRequestInner{Id: v.Id, Start: v.Start, End: v.End, WrappedTask: wt}
	}

	raw, err := json.Marshal(&ws)
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}

func (s *TasksScheduler) UnmarshalJSON(raw []byte) error {
	aa := &innerSequentialSchedule{}

	err := json.Unmarshal(raw, aa)
	if err != nil {
		return err
	}

	adminInterface := reflect.TypeOf((*monitorInterfaces.AdminClient)(nil)).Elem()

	s.Schedule = make(map[string]TaskRequestInfo)
	for k, v := range aa.Schedule {
		t, err := s.marshaller.UnwrapInstance(v.WrappedTask)
		if err != nil {
			return err
		}

		// Marshalling service handlers is mostly non-sense, so
		isAdminClient := reflect.TypeOf(t).Implements(adminInterface)
		if isAdminClient {
			adminClient := t.(monitorInterfaces.AdminClient)
			adminClient.SetAdminHandler(s.adminHandler)
		}

		s.Schedule[k] = TaskRequestInfo{Id: v.Id, Start: v.Start, End: v.End, Task: t.(tasks.Task)}
	}

	return nil
}