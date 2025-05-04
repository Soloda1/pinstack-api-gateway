package user_handler

import (
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	_, err = h.userClient.GetUser(r.Context(), id)
	if err != nil {
		switch err {
		case custom_errors.ErrUserNotFound:
			utils.SendError(w, http.StatusNotFound, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	err = h.userClient.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case custom_errors.ErrUserNotFound:
			utils.SendError(w, http.StatusNotFound, err.Error())
		case custom_errors.ErrOperationNotAllowed:
			utils.SendError(w, http.StatusForbidden, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusNoContent, nil)
}
