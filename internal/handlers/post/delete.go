package post_handler

import (
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Delete godoc
// @Summary Delete a post
// @Description Delete an existing post by ID
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string "Post deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.log.Debug("Missing post id in path params")
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Debug("Invalid post id format", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	err = h.postClient.DeletePost(r.Context(), claims.UserID, id)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrPostNotFound):
			h.log.Debug("delete post failed, not found", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrValidationFailed):
			h.log.Debug("delete post failed, validation failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		case errors.Is(err, custom_errors.ErrForbidden):
			h.log.Debug("delete post failed, forbidden", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
			return
		default:
			h.log.Error("delete post failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
	}
	utils.Send(w, http.StatusOK, nil)
}
