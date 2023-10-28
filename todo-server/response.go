package main

import (
	"encoding/json"
	"github.com/zawlinnnaing/go-clis/to-do/todo"
	"time"
)

type TodoResponse struct {
	Data todo.TaskList `json:"data"`
}

func (resp *TodoResponse) MarshalJSON() ([]byte, error) {
	_resp := struct {
		Data         todo.TaskList `json:"data"`
		Date         int64         `json:"date"`
		TotalResults int           `json:"total_results"`
	}{
		Data:         resp.Data,
		Date:         time.Now().Unix(),
		TotalResults: len(resp.Data),
	}
	return json.Marshal(_resp)
}
