package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/apollobing/todo-app/internal/core/logger"
	core_http_request "github.com/apollobing/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/apollobing/todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask       godoc
// @Summary      Get task
// @Description  Get a specific task by ID
// @Tags         tasks
// @Produce      json
// @Param        id    path      int  true "ID of the task to retrieve"
// @Success      200   {object}  GetTaskResponse "Task retrieved successfully"
// @Failure      400   {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404   {object}  core_http_response.ErrorResponse "Task not found"
// @Failure      500   {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
