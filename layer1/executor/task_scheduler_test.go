package executor

import (
	"fmt"
	"github.com/alicenet/alicenet/constants"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dkg"
	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/alicenet/alicenet/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getTaskScheduler(t *testing.T) (*TasksScheduler, chan tasks.TaskRequest) {
	db := mocks.NewTestDB()
	eth := mocks.NewMockClient()
	adminHandlers := mocks.NewMockAdminHandler()
	txWatcher := transaction.NewWatcher(eth, 12, db, false)
	taskRequestChan := make(chan tasks.TaskRequest, constants.TaskSchedulerBufferSize)
	tasksScheduler, err := NewTasksScheduler(db, eth, adminHandlers, taskRequestChan, txWatcher)
	assert.Nil(t, err)
	return tasksScheduler, taskRequestChan
}

func TestTasksScheduler_Schedule_NilTask(t *testing.T) {
	scheduler, tasksChan := getTaskScheduler(t)
	defer close(tasksChan)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic")
			t.Fatal(r)
		}
	}()

	request := tasks.TaskRequest{
		Task:   nil,
		Action: tasks.Schedule,
	}

	tasksChan <- request

	go func() {
		select {
		case <-time.After(1 * time.Second):
		}
		assert.Emptyf(t, scheduler.Schedule, "Expected zero tasks scheduled")
		scheduler.cancelChan <- true
	}()

	scheduler.eventLoop()
}

func TestTasksScheduler_Schedule_WrongExecutionData(t *testing.T) {

	scheduler, tasksChan := getTaskScheduler(t)
	defer close(tasksChan)

	task := dkg.NewCompletionTask(2, 1)
	request := tasks.TaskRequest{
		Task:   task,
		Action: tasks.Schedule,
	}
	tasksChan <- request

	scheduler.LastHeightSeen = 12
	task = dkg.NewCompletionTask(2, 3)
	request = tasks.TaskRequest{
		Task:   task,
		Action: tasks.Schedule,
	}
	tasksChan <- request

	go func() {
		select {
		case <-time.After(1 * time.Second):
		}
		assert.Emptyf(t, scheduler.Schedule, "Expected zero tasks scheduled")
		scheduler.cancelChan <- true
	}()

	scheduler.eventLoop()
}

func TestTasksScheduler_Schedule_TasksSuccessfully(t *testing.T) {

	scheduler, tasksChan := getTaskScheduler(t)
	defer close(tasksChan)

	scheduler.Start()

	taskCompletionTask := dkg.NewCompletionTask(2, 3)
	request := tasks.TaskRequest{
		Task:   taskCompletionTask,
		Action: tasks.Schedule,
	}
	tasksChan <- request

	taskCompletionTaskSecond := dkg.NewCompletionTask(3, 4)
	request = tasks.TaskRequest{
		Task:   taskCompletionTaskSecond,
		Action: tasks.Schedule,
	}
	tasksChan <- request

	taskRegister := dkg.NewRegisterTask(2, 5)
	request = tasks.TaskRequest{
		Task:   taskRegister,
		Action: tasks.Schedule,
	}
	tasksChan <- request
	select {
	case <-time.After(500 * time.Millisecond):
		assert.Equalf(t, 3, len(scheduler.Schedule), "Expected 3 task scheduled")
	}

	request = tasks.NewKillTaskRequest(&dkg.CompletionTask{})
	tasksChan <- request
	select {
	case <-time.After(500 * time.Millisecond):
		assert.Equalf(t, 1, len(scheduler.Schedule), "Expected 1 task after Completion tasks have been killed")
	}

	request = tasks.NewKillTaskRequest(&dkg.DisputeMissingGPKjTask{})
	tasksChan <- request
	select {
	case <-time.After(500 * time.Millisecond):
		assert.Equalf(t, 1, len(scheduler.Schedule), "There should be 1 tasks left still, due there were no DisputeMissing task scheduled")
	}

	request = tasks.NewKillTaskRequest(&dkg.RegisterTask{})
	tasksChan <- request
	select {
	case <-time.After(500 * time.Millisecond):
		assert.Equalf(t, 0, len(scheduler.Schedule), "All the tasks should have been removed")
	}

	assert.Emptyf(t, scheduler.Schedule, "Expected zero tasks scheduled")
	scheduler.cancelChan <- true
}

//func TestTasksScheduler_Schedule_Success(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	assert.Emptyf(t, s.Schedule, "Expected Schedule map to be empty")
//	err := s.schedule(ctx, task)
//	assert.Nil(t, err)
//	assert.Equal(t, 1, len(s.Schedule))
//
//	for k, _ := range s.Schedule {
//		taskRequest := s.Schedule[k]
//		taskRequest.Start = task.GetStart()
//		taskRequest.End = task.GetEnd()
//	}
//}
//
//func TestTasksScheduler_ProcessTaskResponse_RemoveTaskWithErrNotNil(t *testing.T) {
//
//	s := getTaskScheduler()
//	task := mocks.NewMockITask()
//	ctx, _ := context.WithCancel(context.Background())
//
//	taskResponse := interfaces.TaskResponse{Id: "1", Err: ErrTaskExpired}
//	s.Schedule["1"] = TaskRequestInfo{"First", 1, 1, true, task}
//	assert.NotEmptyf(t, s.Schedule, "Expected one task request scheduled")
//
//	err := s.processTaskResponse(ctx, taskResponse)
//	assert.Nil(t, err)
//	assert.Emptyf(t, s.Schedule, "Expected no tasks")
//}
//
//func TestTasksScheduler_ProcessTaskResponse_RemoveNotScheduledTask(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	taskResponse := interfaces.TaskResponse{Id: "2", Err: ErrTaskExpired}
//	s.Schedule["1"] = TaskRequestInfo{"First", 1, 1, true, task}
//	assert.NotEmptyf(t, s.Schedule, "Expected one task request scheduled")
//
//	err := s.processTaskResponse(ctx, taskResponse)
//	assert.Equal(t, err, ErrNotScheduled)
//	assert.NotEmptyf(t, s.Schedule, "Expected one task to still be in scheduled")
//}
//
//func TestTasksScheduler_ProcessTaskResponse_RemoveTaskWithErrNil(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	taskResponse := interfaces.TaskResponse{Id: "1", Err: nil}
//	s.Schedule["1"] = TaskRequestInfo{"First", 1, 1, true, task}
//	assert.NotEmptyf(t, s.Schedule, "Expected one task request scheduled")
//
//	err := s.processTaskResponse(ctx, taskResponse)
//	assert.Nil(t, err)
//	assert.Emptyf(t, s.Schedule, "Expected no tasks")
//}
//
//func TestTasksScheduler_ProcessTaskResponse_RemoveTaskWithErrNilNotScheduledTask(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	taskResponse := interfaces.TaskResponse{Id: "2", Err: nil}
//	s.Schedule["1"] = TaskRequestInfo{"First", 1, 1, true, task}
//	assert.NotEmptyf(t, s.Schedule, "Expected one task request scheduled")
//
//	err := s.processTaskResponse(ctx, taskResponse)
//	assert.Equal(t, err, ErrNotScheduled)
//	assert.NotEmptyf(t, s.Schedule, "Expected one task to still be in scheduled")
//}
//
//func TestTasksScheduler_StartTask_Success(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 2, false, task},
//		{"Second", 1, 2, false, task},
//	}
//
//	err := s.startTasks(ctx, taskRequestList)
//	assert.Nil(t, err)
//	for _, taskRequest := range taskRequestList {
//		taskRequest.isRunning = true
//		assert.Truef(t, taskRequest.isRunning, "Expecting task to be running")
//	}
//}
//
//func TestTasksScheduler_KillTaskByName(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task1 := mocks.NewMockITask()
//	task2 := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 2, true, task1},
//		{"Second", 1, 2, true, task2},
//	}
//	for _, taskRequest := range taskRequestList {
//		s.Schedule[taskRequest.Id] = taskRequest
//	}
//
//	err := s.killTaskByName(ctx, taskGroupName)
//	assert.Nil(t, err)
//	mockrequire.CalledN(t, task1.Close, 3)
//	mockrequire.CalledN(t, task2.Close, 3)
//}
//
//func TestTasksScheduler_KillTasks(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task1 := mocks.NewMockITask()
//	task2 := mocks.NewMockITask()
//	task3 := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 2, true, task1},
//		{"Second", 1, 2, true, task2},
//		{"Third", 1, 2, true, task3},
//	}
//
//	err := s.killTasks(ctx, taskRequestList)
//	assert.Nil(t, err)
//	killedTaskList := s.findTasksByName(taskGroupName)
//	assert.Emptyf(t, killedTaskList, "Expected no tasks with this name to be running")
//}
//
//func TestTasksScheduler_RemoveUnresponsiveTasks_WithScheduledAndNoScheduledTask(t *testing.T) {
//
//	s := getTaskScheduler()
//	ctx, _ := context.WithCancel(context.Background())
//	task := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 10, true, task},
//		{"Second", 1, 10, true, task},
//		{"Third", 1, 10, true, task},
//	}
//	s.LastHeightSeen = 100
//
//	for _, taskRequest := range taskRequestList {
//		s.Schedule[taskRequest.Id] = taskRequest
//	}
//
//	err := s.removeUnresponsiveTasks(ctx, taskRequestList)
//	assert.Nil(t, err)
//	assert.Emptyf(t, s.Schedule, "Expected no tasks")
//}
//
//func TestTasksScheduler_Purge(t *testing.T) {
//
//	s := getTaskScheduler()
//	task := mocks.NewMockITask()
//	ctx, _ := context.WithCancel(context.Background())
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 2, true, task},
//		{"Second", 1, 2, true, task},
//		{"Third", 1, 2, true, task},
//	}
//
//	for _, taskRequest := range taskRequestList {
//		err := s.schedule(ctx, taskRequest.Task)
//		assert.Nil(t, err)
//	}
//	s.purge()
//	assert.Emptyf(t, s.Schedule, "Expected no tasks")
//}
//
//func TestTasksScheduler_FindTasks(t *testing.T) {
//
//	s := getTaskScheduler()
//	task := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 10, true, task},
//		{"Second", 1, 90, true, task},
//		{"Third", 0, 0, false, task},
//		{"Fourth", 10, 0, false, task},
//		{"Fifth", 1, 200, false, task},
//	}
//	s.LastHeightSeen = 100
//
//	for _, taskRequest := range taskRequestList {
//		s.Schedule[taskRequest.Id] = taskRequest
//	}
//
//	toStart, expired, unresponsive := s.findTasks()
//	assert.Equalf(t, 3, len(toStart), "Expected 3 tasks to start")
//	assert.Equalf(t, 1, len(expired), "Expected 1 task expired")
//	assert.Equalf(t, 1, len(unresponsive), "Expected 1 task unresponsive")
//}
//
//func TestTasksScheduler_FindTasksByName(t *testing.T) {
//
//	s := getTaskScheduler()
//	taskWithGroupName := mocks.NewMockITask()
//	taskWithoutGroupName := mocks.NewMockITask()
//
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 2, true, taskWithGroupName},
//		{"Second", 1, 2, true, taskWithGroupName},
//		{"Third", 1, 2, true, taskWithoutGroupName},
//	}
//	for _, taskRequest := range taskRequestList {
//		s.Schedule[taskRequest.Id] = taskRequest
//	}
//
//	tasksByName := s.findTasksByName(taskGroupName)
//	assert.Equalf(t, 2, len(tasksByName), "Expected 2 tasks with this name")
//}
//
//func TestTasksScheduler_Remove(t *testing.T) {
//
//	s := getTaskScheduler()
//	task := mocks.NewMockITask()
//
//	s.Schedule["First"] = TaskRequestInfo{"First", 1, 2, true, task}
//	err := s.remove("First")
//	assert.Nil(t, err)
//	assert.Emptyf(t, s.Schedule, "Expected task to be deleted")
//
//	err = s.remove("Second")
//	assert.NotNil(t, err)
//	assert.Equal(t, err, ErrNotScheduled)
//}

//func TestTasksScheduler_PersistState_Success(t *testing.T) {
//
//	s := getTaskScheduler()
//
//	s.LastHeightSeen = 0
//	err := s.persistState()
//	assert.Nil(t, err)
//
//	err = s.loadState()
//	assert.Nil(t, err)
//
//	lastHeightSeenBeforeAfter := s.LastHeightSeen
//	s.LastHeightSeen = lastHeightSeen
//	err = s.persistState()
//	err = s.loadState()
//	lastHeightSeenAfter := s.LastHeightSeen
//	assert.Nil(t, err)
//	assert.NotEqualf(t, lastHeightSeenBeforeAfter, lastHeightSeenAfter, "Expected TaskScheduler to be different")
//}
//
//func TestTasksScheduler_LoadState_MissingKey(t *testing.T) {
//
//	s := getTaskScheduler()
//
//	err := s.loadState()
//	assert.NotNil(t, err)
//}
//
//func TestTasksScheduler_Start_EmptySchedule(t *testing.T) {
//
//	s := getTaskScheduler()
//
//	err := s.persistState()
//	assert.Nil(t, err)
//
//	assert.Emptyf(t, s.Schedule, "Scheduled map expected to be empty")
//	err = s.Start()
//	assert.Nil(t, err)
//	assert.Emptyf(t, s.Schedule, "Scheduled map expected to still be empty")
//}

//func TestTasksScheduler_EventLoop_PurgeToEmptyTheSchedulerMap(t *testing.T) {
//
//	db := mocks.NewTestDB()
//	eth := mocks.NewMockNetwork()
//	adminHandlers := mocks.NewMockIAdminHandler()
//	txWatcher := transaction.NewWatcher(eth, 12, db, false)
//	taskRequestChan := make(chan interfaces.Task, 100)
//	s := NewTasksScheduler(db, eth, adminHandlers, taskRequestChan, nil, txWatcher)
//	task := mocks.NewMockITask()
//
//	s.Schedule["1"] = TaskRequestInfo{"First", 1, 1, true, task}
//	s.Schedule["2"] = TaskRequestInfo{"Second", 2, 2, true, task}
//	assert.Equal(t, 2, len(s.Schedule))
//	s.purge()
//
//	assert.Equal(t, 0, len(s.Schedule))
//}
//
//func TestTasksScheduler_Schedule_ScheduleTask(t *testing.T) {
//
//	db := mocks.NewTestDB()
//	eth := mocks.NewMockNetwork()
//	adminHandlers := mocks.NewMockIAdminHandler()
//	taskRequestChan := make(chan interfaces.Task, 100)
//	txWatcher := transaction.NewWatcher(eth, 12, db, false)
//	s := NewTasksScheduler(db, eth, adminHandlers, taskRequestChan, nil, txWatcher)
//
//	s.LastHeightSeen = lastHeightSeen
//	task := &dkgtasks.CompletionTask{BaseTask: objects.NewBaseTask(taskGroupName, 10, 20, false, nil)}
//
//	err := s.persistState()
//	assert.Nil(t, err)
//	s2 := s
//	assert.Equal(t, s2, s)
//
//	taskRequestChan <- task
//	go func() {
//		time.Sleep(10 * time.Second)
//		s.cancelChan <- true
//		s.logger.Debugf("s.Schedule: %v", s.Schedule)
//		assert.NotEmpty(t, s.Schedule)
//	}()
//
//	s.eventLoop()
//}
//
//func TestTasksScheduler_EventLoop_Workflow(t *testing.T) {
//
//	db := mocks.NewTestDB()
//	eth := mocks.NewMockNetwork()
//	adminHandlers := mocks.NewMockIAdminHandler()
//	taskRequestChan := make(chan interfaces.Task, 100)
//	taskKillChan := make(chan string, 1)
//	txWatcher := transaction.NewWatcher(eth, 12, db, false)
//	s := NewTasksScheduler(db, eth, adminHandlers, taskRequestChan, taskKillChan, txWatcher)
//	s.LastHeightSeen = lastHeightSeen
//
//	// Create a valid task
//	task := &dkgtasks.CompletionTask{BaseTask: objects.NewBaseTask(taskGroupName, 10, 20, false, nil)}
//
//	// Send task to request channel
//	taskRequestList := []TaskRequestInfo{
//		{"First", 1, 10, false, task},
//		{"Second", 1, 90, false, task},
//		{"Third", 0, 0, false, task},
//		{"Fourth", 10, 0, false, task},
//		{"Fifth", 1, 200, false, task},
//	}
//	for _, taskRequest := range taskRequestList {
//		taskRequestChan <- taskRequest.Task
//	}
//
//	// Initial State
//	err := s.loadState()
//	assert.NotNil(t, err)
//
//	// Execute after processing time
//	eth.GetCurrentHeightFunc.PushReturn(15, nil)
//
//	// Close scheduler
//	go func() {
//		time.Sleep(7 * time.Second)
//		s.cancelChan <- true
//	}()
//
//	// Expect event loop to clean the tasks
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		s.eventLoop()
//	}()
//	wg.Wait()
//
//	mockrequire.CalledN(t, eth.GetCurrentHeightFunc, 1)
//	assert.Emptyf(t, s.Schedule, "Expected all the task to have been removed")
//}
