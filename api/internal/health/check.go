package health

type Status struct {
	State string
}

type CheckHealth struct{}

func (CheckHealth) Execute() Status {
	return Status{State: "ok"}
}
