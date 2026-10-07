package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/apollobing/todo-app/internal/core/domain"
	core_logger "github.com/apollobing/todo-app/internal/core/logger"
	core_http_request "github.com/apollobing/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/apollobing/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/apollobing/todo-app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"       swaggertype:"string" example:"To walk the dog"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}

		if r.Description.Set {
			if r.Description.Value != nil {
				descriptionLen := len([]rune(*r.Description.Value))
				if descriptionLen < 1 || descriptionLen > 1000 {
					return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
				}
			}
		}

		if r.Completed.Set {
			if r.Completed.Value == nil {
				return fmt.Errorf("`Completed` can't be NULL")
			}
		}
	}

	return nil
}

type PatchUserResponse TaskDTOResponse

// PatchTask     godoc
// @Summary      Update task
// @Description  Update information of an existing task in the system
// @Description  ### Field update logic (Three-state logic):
// @Description  1. **Field omitted**: `description` is ignored, DB value remains unchanged
// @Description  2. **Explicit value provided**: `"description": "Go for a walk with Milo at 06:30 AM"` - sets the new description in the DB
// @Description  3. **null provided**: `"description": null` - clears the field in the DB (set to NULL)
// @Description  Constraints: `title` and `completed` cannot be set to null
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "ID of the task to update"
// @Param        request  body      PatchTaskRequest  true  "PatchTask request body"
// @Success      200      {object}  PatchUserResponse                "Task updated successfully"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404      {object}  core_http_response.ErrorResponse "Task not found"
// @Failure      409      {object}  core_http_response.ErrorResponse "Conflict"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
