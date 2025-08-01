package notification_handler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
)

type ReadAllUserNotificationsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ReadAllUserNotifications godoc
// @Summary Mark all user notifications as read
// @Description Mark all notifications of a user as read
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ReadAllUserNotificationsResponse "All notifications marked as read"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/read-all [put]
func (h *NotificationHandler) ReadAllUserNotifications(w http.ResponseWriter, r *http.Request) {
	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	err = h.notificationClient.ReadAllUserNotifications(r.Context(), claims.UserID)
	if err != nil {
		h.log.Error("Failed to mark all notifications as read", slog.Int64("user_id", claims.UserID), slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrNotificationNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	response := ReadAllUserNotificationsResponse{
		Success: true,
		Message: "All notifications marked as read",
	}
	utils.Send(w, http.StatusOK, response)
}
