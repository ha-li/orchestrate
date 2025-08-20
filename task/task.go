package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type (
	Task struct {
		ID            uuid.UUID // uuid are 128 bits, in practice unique, but it is possible to generate to identical uuids
		Name          string
		State         State
		Image         string
		Memory        int
		Disk          int
		ExposedPorts  nat.PortSet
		PortBindings  map[string]string
		RestartPolicy string
		StartTime     time.Time
		FinishTime    time.Time
	}

	TaskEvent struct {
		ID        uuid.UUID
		State     State
		Timestamp time.Time
		Task      Task
	}
)
