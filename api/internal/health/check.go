package health

type Status struct {
	State string
}

type CheckHealth struct{}

func NewCheckHealth() CheckHealth {
	return CheckHealth{}
}

func (CheckHealth) Execute() Status {
	return Status{State: "ok"}
}
