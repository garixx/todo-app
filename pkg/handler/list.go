package handler

import (
	"github.com/gin-gonic/gin"
	"github/garixx/todo-app"
	"net/http"
	"strconv"
)

// @Summary Create Todo list
// @Security ApiKeyAuth
// @Tags api
// @Description create todo list
// @ID create-todo-list
// @Accept json
// @Produce json
// @Param input body todo.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400, 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, "invalid id param")
		return
	}

	todoList, err := h.services.TodoList.GetById(userId, listId)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, todoList)
}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
