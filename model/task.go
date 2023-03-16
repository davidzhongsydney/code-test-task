package model

import timestamp "github.com/golang/protobuf/ptypes/timestamp"

type Task struct {
	TaskID    uint64               `json:"TaskID,omitempty"`
	Name      string               `json:"Name,omitempty"`
	Content   string               `json:"Content,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"CreatedAt,omitempty"`
	UpdatedAt *timestamp.Timestamp `json:"UpdatedAt,omitempty"`
	DeletedAt *timestamp.Timestamp `json:"DeletedAt,omitempty"`
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

func (x *Task) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Task) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Task) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}
