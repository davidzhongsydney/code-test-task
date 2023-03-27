package model

import (
	"time"
)

type T_Task struct {
	Task
	T_Internal
}

type Task struct {
	TaskID  uint64 `json:"taskID,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type T_Internal struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (x *Task) GetTaskID() uint64 {
	if x != nil {
		return x.TaskID
	}
	return 0
}

func (x *Task) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Task) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *T_Internal) GetCreatedAt() *time.Time {
	if x != nil {
		return x.CreatedAt
	}
	return &time.Time{}
}

func (x *T_Internal) GetUpdatedAt() *time.Time {
	if x != nil {
		return x.UpdatedAt
	}
	return &time.Time{}
}

func (x *T_Internal) GetDeletedAt() *time.Time {
	if x != nil {
		return x.DeletedAt
	}
	return &time.Time{}
}
