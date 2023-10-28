package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		replyError(writer, request, http.StatusNotFound, "Route not found")
		return
	}
	content := "Hello from root route"
	replyTextContent(writer, request, http.StatusOK, content)
}

func todoRouter(todoFile string, locker sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.TaskList{}
		locker.Lock()
		defer locker.Unlock()
		err := list.Load(todoFile)
		if err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			case http.MethodPost:
				addHandler(w, r, list, todoFile)
			default:
				replyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
			}
			return
		}
		id, err := validateID(r.URL.Path, list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, ErrNotFound.Error())
				return
			}
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, id)
		case http.MethodPatch:
			patchOneHandler(w, r, list, id, todoFile)
		case http.MethodDelete:
			deleteOneHandler(w, r, list, id, todoFile)
		default:
			replyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
		}

	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.TaskList) {
	resp := &TodoResponse{Data: *list}

	replyJSONContent(w, r, http.StatusOK, resp)
}

func addHandler(w http.ResponseWriter, r *http.Request, list *todo.TaskList, todoFile string) {
	item := struct {
		Task string `json:"task"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		replyError(w, r, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %s", err.Error()))
		return
	}

	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, fmt.Sprintf("Save file error: %s", err.Error()))
		return
	}

	replyTextContent(w, r, http.StatusCreated, "")
}

func validateID(id string, list *todo.TaskList) (int, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("%w: InvalidData %s", ErrInvalidData, err.Error())
	}
	if idInt < 0 {
		return 0, fmt.Errorf("%w: ID less than 1", ErrInvalidData)
	}
	if idInt >= len(*list) {
		return 0, fmt.Errorf("%w: ID greater than list length(%d)", ErrNotFound, len(*list))
	}
	return idInt, nil
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.TaskList, id int) {
	resp := &TodoResponse{
		Data: (*list)[id : id+1],
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}
func deleteOneHandler(w http.ResponseWriter, r *http.Request, list *todo.TaskList, id int, todoFile string) {
	if err := list.Delete(id); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}
func patchOneHandler(w http.ResponseWriter, r *http.Request, list *todo.TaskList, id int, todoFile string) {
	query := r.URL.Query()

	if _, ok := query["complete"]; !ok {
		replyError(w, r, http.StatusBadRequest, "Missing complete query")
		return
	}

	if err := list.Complete(id); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	replyTextContent(w, r, http.StatusNoContent, "")
}
