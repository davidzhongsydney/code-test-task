package encoder

type ErrorMessage string

const (
	TASK_NOT_EXIST        ErrorMessage = "task does not exist"
	TASK_DELETED          ErrorMessage = "task has been deleted"
	TASK_ID_NOT_SPECIFIED ErrorMessage = "task id not specified"
)
