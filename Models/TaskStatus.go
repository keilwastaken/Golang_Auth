package Models

type TaskStatus string

const (
	Created   TaskStatus = "created"
	Pause     TaskStatus = "pause"
	Complete  TaskStatus = "complete"
	Cancelled TaskStatus = "cancelled"
)
