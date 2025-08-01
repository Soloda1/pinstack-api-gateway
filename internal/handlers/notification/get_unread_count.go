package notification_handler

import (
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUnreadCountResponse struct {
	Count int32 `json:"count"`
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the count of unread notifications for a user
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} GetUnreadCountResponse "Unread notification count"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/unread-count/{user_id} [get]
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Debug("Failed to parse user ID", slog.String("user_id", userIDStr), slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	if claims.UserID != userID {
		h.log.Debug(
			"User not authorized to get unread notifications count for other users",
			slog.Int64("request_user_id", userID),
			slog.Int64("authenticated_user_id", claims.UserID),
		)
		utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
		return
	}

	count, err := h.notificationClient.GetUnreadCount(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get unread notification count", slog.Int64("user_id", userID), slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrNotificationAccessDenied):
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	response := GetUnreadCountResponse{
		Count: count,
	}
	utils.Send(w, http.StatusOK, response)
}
