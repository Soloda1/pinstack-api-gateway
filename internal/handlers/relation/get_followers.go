package relation_handler

import (
	"errors"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetFollowersResponse struct {
	Followers []*models.RelationUser `json:"followers"`
	Total     int64                  `json:"total"`
	Page      int32                  `json:"page"`
	Limit     int32                  `json:"limit"`
}

// GetFollowers godoc
// @Summary Get user followers
// @Description Get list of user followers by user ID
// @Tags relation
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Limit for pagination" default(20)
// @Param page query int false "Page number" default(1)
// @Success 200 {object} GetFollowersResponse "Followers retrieved successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /relation/{user_id}/followers [get]
func (h *RelationHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Debug("Failed to parse user ID", slog.String("user_id", userIDStr), slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	if userID <= 0 {
		h.log.Debug("Invalid user ID", slog.Int64("user_id", userID))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit := int32(20) // default
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 && l <= 100 {
			limit = int32(l)
		}
	}

	page := int32(1) // default
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 32); err == nil && p > 0 {
			page = int32(p)
		}
	}

	followers, total, err := h.relationClient.GetFollowers(r.Context(), userID, limit, page)
	if err != nil {
		h.log.Error("Failed to get followers", slog.Int64("user_id", userID), slog.String("error", err.Error()))

		switch {
		case errors.Is(err, custom_errors.ErrValidationFailed):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrDatabaseQuery):
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrDatabaseQuery.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	response := GetFollowersResponse{
		Followers: followers,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}
	utils.Send(w, http.StatusOK, response)
}
