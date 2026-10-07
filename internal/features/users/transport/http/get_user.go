package users_transport_http

import (
	"net/http"

	core_logger "github.com/apollobing/todo-app/internal/core/logger"
	core_http_response "github.com/apollobing/todo-app/internal/core/transport/http/response"
	core_http_request "github.com/apollobing/todo-app/internal/core/transport/http/request"
)

type GetUserResponse UserDTOResponse

// GetUser       godoc
// @Summary      Get user
// @Description  Get a specific user by ID
// @Tags         users
// @Produce      json
// @Param        id path int true                              "ID of the user to retrieve"
// @Success      200 {object} GetUserResponse                  "User retrieved successfully"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "User not found"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
