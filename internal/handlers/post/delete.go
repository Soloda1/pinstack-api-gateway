package post_handler

import (
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
	"strconv"
)

// Delete godoc
// @Summary Delete a post
// @Description Delete an existing post by ID
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id query string true "Post ID"
// @Success 200 {object} map[string]string "Post deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [delete]
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.log.Debug("Missing post id in query params")
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Debug("Invalid post id format", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}
	claimsRaw := r.Context().Value("claims")
	claims, ok := claimsRaw.(*middlewares.Claims)
	if !ok || claims == nil {
		h.log.Error("invalid token claims")
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidToken.Error())
		return
	}
	err = h.postClient.DeletePost(r.Context(), claims.UserID, id)
	if err != nil {
		h.log.Error("delete post failed", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
		return
	}
	utils.Send(w, http.StatusOK, nil)
}
