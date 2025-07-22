package job

type JobStatus int

const (
	Pending JobStatus = iota
	InProgress
	Done
	Failed
)

var jobStatusName = map[JobStatus]string{
	Pending:    "pending",
	InProgress: "in_progress",
	Done:       "done",
	Failed:     "failed",
}

var jobStatusValue = map[string]JobStatus{
	"pending":     Pending,
	"in_progress": InProgress,
	"done":        Done,
	"failed":      Failed,
}

func (s JobStatus) String() string {
	if name, ok := jobStatusName[s]; ok {
		return name
	}
	return "unknown"
}

func ParseStatus(s string) (JobStatus, bool) {
	status, ok := jobStatusValue[s]
	return status, ok
}
