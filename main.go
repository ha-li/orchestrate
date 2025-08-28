package main

import (
	"fmt"
	"log/slog"
	"orchestrate/app/config"
	"orchestrate/task"
	"orchestrate/worker"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	/*  // previous main
	t := task.Task{
			ID:     uuid.New(),
			Name:   "Task-1",
			State:  task.Pending,
			Image:  "Image-1",
			Memory: 1024,
			Disk:   1,
		}

		te := task.TaskEvent{
			ID:        uuid.New(),
			State:     task.Pending,
			Timestamp: time.Now(),
			Task:      t,
		}

		fmt.Printf("task: %v\n", t)
		fmt.Printf("task event: %v\n", te)

		w := worker.Worker{
			Name:  "worker-1",
			Queue: *queue.New(),
			Db:    make(map[uuid.UUID]*task.Task),
		}

		fmt.Printf("worker: %v\n", w)
		w.CollectStats()
		w.RunTask()
		w.StartTask()
		w.StopTask()

		m := manager.Manager{
			Pending: *queue.New(),
			TaskDb:  make(map[string][]*task.Task),
			EventDb: make(map[string][]*task.TaskEvent),
			Workers: []string{w.Name},
		}

		fmt.Printf("manager: %v\n", m)
		m.SelectWorker()
		m.UpdateTasks()
		m.SendWork()

		n := node.Node{
			Name:   "Node-1",
			Ip:     "192.168.1.1",
			Cores:  4,
			Memory: 1024,
			Disk:   25,
			Role:   "worker",
		}
		fmt.Printf("node: %v\n", n)

		fmt.Printf("create a test contaienr\n")
		dockerTask, createResult := createContainer()
		if createResult.Error != nil {
			fmt.Printf("createContainer error: %v\n", createResult.Error)
			os.Exit(1)
		}

		fmt.Printf("createResult: %v\n", createResult)
		time.Sleep(time.Second * 5)
		fmt.Printf("stopping container %s\n", createResult.ContainerId)
		_ = stopContainer(dockerTask, createResult.ContainerId)

	*/
	// chapter 4 stuff
	/* appConfig := config.New()
	err := appConfig.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	// the first task
	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Printf("starting task: %v\n", t)
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
	t.ContainerID = result.ContainerId
	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 30)

	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	} */

	// chapter 5
	host := os.Getenv("localhost")
	port, _ := strconv.Atoi("5555")

	slog.Info("Starting cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	api := worker.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w) // a goroutine is just a new thread
	// doesn't work on mac since it doesn't have /proc
	//go w.CollectStats()
	api.Start()
}

// In golang, a general guideline of when to use exported functions vs structs and methods is
// if working with a stateless operation, use functions.
// If the operation requires encapsulate state, or a clear lifecycle, then it is appropriate
// to use a struct with methods. You can use fields to hold state/data, and enforce lifecycle.
func runTasks(w *worker.Worker) {

	// in infinite loop, with sleeps
	for {
		// if there is something in the queue, we have the worker run a task
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				slog.Error("Error running tasks", "error", result.Error)
			}
		} else {
			// otherwise we print to console no task
			slog.Info("No tasks to process currently")
		}

		// sleep until the next check
		slog.Info("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

func createContainer(appConfig *config.AppConfig) (*task.Docker, *task.DockerResult) {

	c := task.Config{
		Name:  appConfig.GetDockerConfig().GetName(),
		Image: appConfig.GetDockerConfig().GetImage(),
		Env: []string{
			"POSTGRES_USER=" + appConfig.GetDatabaseConfig().GetUser(),
			"POSTGRES_PASSWORD=" + appConfig.GetDatabaseConfig().GetPassword(),
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, &result
	}

	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return &result
	}

	fmt.Printf("Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}
