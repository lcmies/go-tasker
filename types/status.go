package types

type Status int

const (
	PENDING Status = iota
	RUNNING
	CANCELLED
	DONE
	ERROR
)

func (s Status) Str() string {
	switch s {
	case PENDING:
		return "PENDING"
	case RUNNING:
		return "RUNNING"
	case CANCELLED:
		return "CANCELLED"
	case DONE:
		return "DONE"
	case ERROR:
		return "ERROR"
	}
	return ""
}
