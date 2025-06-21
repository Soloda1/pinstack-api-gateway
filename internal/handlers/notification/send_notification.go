package notification_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SendNotificationRequest struct {
	UserID  int64           `json:"user_id" validate:"required,gt=0"`
	Type    string          `json:"type" validate:"required"`
	Payload json.RawMessage `json:"payload"`
}

type SendNotificationResponse struct {
	NotificationID int64  `json:"notification_id"`
	Message        string `json:"message"`
}

// SendNotification godoc
// @Summary Send notification
// @Description Send notification to a user
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendNotificationRequest true "Send notification request"
// @Success 200 {object} SendNotificationResponse "Notification sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 409 {object} map[string]string "Notification limit exceeded"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/send [post]
func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode send notification request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.StructPartial(req, "UserID", "Type"); err != nil {
		h.log.Debug("Failed to validate send notification request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	notificationID, err := h.notificationClient.SendNotification(r.Context(), req.UserID, req.Type, req.Payload)
	if err != nil {
		h.log.Error("send notification failed", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
				return
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
				return
			case codes.ResourceExhausted:
				utils.SendError(w, http.StatusTooManyRequests, custom_errors.ErrNotificationLimitExceeded.Error())
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, custom_errors.ErrInsufficientRights.Error())
				return
			case codes.AlreadyExists:
				utils.SendError(w, http.StatusConflict, "Notification already exists")
				return
			case codes.Unavailable:
				utils.SendError(w, http.StatusServiceUnavailable, custom_errors.ErrExternalServiceUnavailable.Error())
				return
			case codes.DeadlineExceeded:
				utils.SendError(w, http.StatusGatewayTimeout, custom_errors.ErrExternalServiceTimeout.Error())
				return
			case codes.Unimplemented:
				utils.SendError(w, http.StatusNotImplemented, "Notification type not implemented")
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrNotificationInvalidType):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrNotificationInvalidType.Error())
		case errors.Is(err, custom_errors.ErrNotificationInvalidPayload):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrNotificationInvalidPayload.Error())
		case errors.Is(err, custom_errors.ErrNotificationLimitExceeded):
			utils.SendError(w, http.StatusConflict, custom_errors.ErrNotificationLimitExceeded.Error())
		case errors.Is(err, custom_errors.ErrNotificationAccessDenied):
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := SendNotificationResponse{
		NotificationID: notificationID,
		Message:        "Notification sent successfully",
	}
	utils.Send(w, http.StatusOK, response)
}
