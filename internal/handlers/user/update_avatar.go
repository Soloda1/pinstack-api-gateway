package user_handler

import (
	"encoding/json"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	err = h.userClient.UpdateAvatar(r.Context(), id, req.AvatarURL)
	if err != nil {
		switch err {
		case custom_errors.ErrUserNotFound:
			utils.SendError(w, http.StatusNotFound, err.Error())
		case custom_errors.ErrValidationFailed:
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case custom_errors.ErrOperationNotAllowed:
			utils.SendError(w, http.StatusForbidden, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, nil)
}
