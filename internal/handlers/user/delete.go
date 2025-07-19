package user_handler

import (
	"errors"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete a user by their ID
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {object} nil "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Operation not allowed"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	_, err = h.userClient.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	err = h.userClient.DeleteUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, custom_errors.ErrOperationNotAllowed):
			utils.SendError(w, http.StatusForbidden, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusNoContent, nil)
}
