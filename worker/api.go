package worker

import (
	"github.com/go-chi/chi"
)

// The API stack
//    Worker        The worker
//    Handlers      provided by chi
//    Router        provided by chi
//    Routes        provided by chi
//    HTTP Server   is the lowest level, provided by http
//     ^  |
// req |  |
//     |  |
//     |  v  response
//    Manager

//   the routes
//   GET     /tasks             get a list of all task, http response 200 (success)
//   POST    /tasks             create a task,                        201 (resource was created, no content)
//   DELETE  /tasks/{taskID}    stop the task id by taskID            204 (no content)

type (
	ErrResponse struct {
		HTTPStatusCode int
		Message        string
	}

	Api struct {
		Address string
		Port    int
		Worker  *Worker
		Router  *chi.Mux // a mux is a multiplexer, synonymous with request router, matches the
		// incoming request path with route to pass the request to the correct handler
	}
)
