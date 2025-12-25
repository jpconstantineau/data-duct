package pipeline

type State string

const (
	StateSucceeded State = "succeeded"
	StateCancelled State = "cancelled"
	StateFailed    State = "failed"
)

type Result interface {
	State() State
	Err() error
}

type Succeeded struct{}

func (Succeeded) State() State { return StateSucceeded }
func (Succeeded) Err() error   { return nil }

type Cancelled struct{ Cause error }

func (c Cancelled) State() State { return StateCancelled }
func (c Cancelled) Err() error   { return c.Cause }

type Failed struct{ Cause error }

func (f Failed) State() State { return StateFailed }
func (f Failed) Err() error   { return f.Cause }
