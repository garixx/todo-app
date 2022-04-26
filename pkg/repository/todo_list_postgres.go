package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github/garixx/todo-app"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (t *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != err {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (t *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT t1.id, t1.title, t1.description FROM %s t1 INNER JOIN %s u1 ON t1.id = u1.list_id WHERE u1.user_id = $1", todoListsTable, usersListsTable)
	err := t.db.Select(&lists, query, userId)
	return lists, err
}

func (t *TodoListPostgres) GetById(userId int, listId int) (todo.TodoList, error) {
	var todoList todo.TodoList

	query := fmt.Sprintf("SELECT t1.id, t1.title, t1.description FROM %s t1 INNER JOIN %s u1 ON t1.id = u1.list_id WHERE u1.user_id = $1 AND u1.list_id = $2", todoListsTable, usersListsTable)
	err := t.db.Get(&todoList, query, userId, listId)
	return todoList, err
}
