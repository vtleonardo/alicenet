package tasks

import "github.com/google/uuid"

// Action is an enumeration indicating the actions that the scheduler
// can do with a task during a request:
type Action int

// The possible actions that the scheduler can do with a task during a request:
// * Kill          - To kill/prune a task type immediately
// * Schedule      - To schedule a new task
const (
	Kill Action = iota
	Schedule
)

func (action Action) String() string {
	return [...]string{
		"Kill",
		"Schedule",
	}[action]
}

type Response struct {
	Id  string
	Err error
}

type Request struct {
	Id     string
	Action Action
	Task   Task
}

func NewRequestScheduleTask(task Task) Request {
	return Request{Action: Schedule, Task: task, Id: uuid.New().String()}
}

func NewRequestKillTaskByType(task Task) Request {
	return Request{Action: Kill, Task: task, Id: ""}
}

func NewRequestKillTaskById(id string) Request {
	return Request{Action: Kill, Task: nil, Id: id}
}
