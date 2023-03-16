package model

import (
	"time"
)

type Task struct {
	TaskID    uint64    `json:"TaskID,omitempty"`
	Name      string    `json:"Name,omitempty"`
	Content   string    `json:"Content,omitempty"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty"`
	DeletedAt time.Time `json:"DeletedAt,omitempty"`
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

func (x *Task) GetCreatedAt() time.Time {
	if x != nil {
		return x.CreatedAt
	}
	return time.Time{}
}

func (x *Task) GetUpdatedAt() time.Time {
	if x != nil {
		return x.UpdatedAt
	}
	return time.Time{}
}

func (x *Task) GetDeletedAt() time.Time {
	if x != nil {
		return x.DeletedAt
	}
	return time.Time{}
}
