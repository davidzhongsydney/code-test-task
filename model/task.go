package model

import timestamp "github.com/golang/protobuf/ptypes/timestamp"

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
	CreatedAt *timestamp.Timestamp `json:"createdAt,omitempty"`
	UpdatedAt *timestamp.Timestamp `json:"updatedAt,omitempty"`
	DeletedAt *timestamp.Timestamp `json:"deletedAt,omitempty"`
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

func (x *T_Internal) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *T_Internal) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *T_Internal) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}
