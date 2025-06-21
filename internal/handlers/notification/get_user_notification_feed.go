package notification_handler

import (
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserNotificationFeedSwaggerResponse is the Swagger response structure for GetUserNotificationFeed
type GetUserNotificationFeedSwaggerResponse struct {
	Notifications []*models.NotificationSwagger `json:"notifications"`
	Total         int                           `json:"total"`
	Page          int                           `json:"page"`
	Limit         int                           `json:"limit"`
	TotalPages    int                           `json:"total_pages"`
}

type GetUserNotificationFeedResponse struct {
	Notifications []*models.Notification `json:"notifications"`
	Total         int                    `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// GetUserNotificationFeed godoc
// @Summary Get user notification feed
// @Description Get paginated list of notifications for a user
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} GetUserNotificationFeedSwaggerResponse "User notification feed"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/feed/{user_id} [get]
func (h *NotificationHandler) GetUserNotificationFeed(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Debug("Failed to parse user ID", slog.String("user_id", userIDStr), slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var page, limit int32 = 1, 10

	if pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			h.log.Debug("Invalid page parameter", slog.String("page", pageStr))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
			return
		}
		page = int32(pageInt)
	}

	if limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 || limitInt > 100 {
			h.log.Debug("Invalid limit parameter", slog.String("limit", limitStr))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
			return
		}
		limit = int32(limitInt)
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	if claims.UserID != userID {
		h.log.Debug("User not authorized to view this feed",
			slog.Int64("request_user_id", userID),
			slog.Int64("authenticated_user_id", claims.UserID))
		utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
		return
	}

	notifications, total, err := h.notificationClient.GetUserNotificationFeed(r.Context(), userID, page, limit)
	if err != nil {
		h.log.Error("Failed to get user notification feed",
			slog.Int64("user_id", userID),
			slog.Int("page", int(page)),
			slog.Int("limit", int(limit)),
			slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
				return
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, custom_errors.ErrInsufficientRights.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		return
	}

	var totalPages int32
	if total%limit == 0 {
		totalPages = total / limit
	} else {
		totalPages = (total / limit) + 1
	}

	response := GetUserNotificationFeedResponse{
		Notifications: notifications,
		Total:         int(total),
		Page:          int(page),
		Limit:         int(limit),
		TotalPages:    int(totalPages),
	}

	utils.Send(w, http.StatusOK, response)
}
