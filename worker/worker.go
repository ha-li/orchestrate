package worker

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"orchestrate/task"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Db        map[uuid.UUID]*task.Task // mapping of task uuid to the task
	Queue     queue.Queue
	TaskCount int
	Stats     *Stats
}

func (w *Worker) GetTasks() []*task.Task {
	tasks := make([]*task.Task, len(w.Db))
	for _, task := range w.Db {
		tasks = append(tasks, task)
	}
	return tasks
}

func (w *Worker) CollectStats() {
	for {
		slog.Info("Collecting stats")
		w.Stats = GetStats()
		w.Stats.TaskCount = w.TaskCount
		time.Sleep(15 * time.Second)
	}
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) RunTask() task.DockerResult {
	fmt.Println("I will run task")
	t := w.Queue.Dequeue() // returns any
	if t == nil {
		log.Println("No task to run")
		return task.DockerResult{Error: nil}
	}

	taskQueued := t.(task.Task) // unwrap any to its true type

	taskPersisted := w.Db[taskQueued.ID] // retrieve the task from the db
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}

	// check if the state transition is valid
	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("Invalid task state")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}
	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	fmt.Println("I will start a task")
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Err running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	c := task.NewConfig(&t)
	d := task.NewDocker(c)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", t.ContainerID, result.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v\n", t.ContainerID, t.ID)

	return result
}
