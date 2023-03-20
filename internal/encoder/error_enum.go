package encoder

type ErrorMessage string

const (
	TASK_NOT_EXIST        ErrorMessage = "task does not exist"
	TASK_DELETED          ErrorMessage = "task has been logically deleted"
	TASK_ID_NOT_SPECIFIED ErrorMessage = "task id not specified"
	TASK_CREATION_ERROR   ErrorMessage = "task is failed to be created"
	TASK_DATABASE_TIMEOUT ErrorMessage = "task database timeout"
)
