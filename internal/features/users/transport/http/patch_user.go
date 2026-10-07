package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apollobing/todo-app/internal/core/domain"
	core_logger "github.com/apollobing/todo-app/internal/core/logger"
	core_http_request "github.com/apollobing/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/apollobing/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/apollobing/todo-app/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+491234567890"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must starts with '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser     godoc
// @Summary      Update user
// @Description  Update information of an existing user in the system
// @Description  ### Field update logic (Three-state logic):
// @Description  1. **Field omitted**: `phone_number` is ignored, DB value remains unchanged
// @Description  2. **Explicit value provided**: `"phone_number": "+491234567890"` - sets the new phone number in the DB
// @Description  3. **null provided**: `"phone_number": null` - clears the field in the DB (set to NULL)
// @Description  Constraints: `full_name` cannot be set to null
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path int              true            "ID of the user to update"
// @Param        request body PatchUserRequest true            "PatchUser request body"
// @Success      200 {object} PatchUserResponse                "User updated successfully"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "User not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
