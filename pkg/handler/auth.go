package handler

import (
	"github.com/gin-gonic/gin"
	"github/garixx/todo-app"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400, 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var input todo.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) signIn(ctx *gin.Context) {
	var input SignInInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
