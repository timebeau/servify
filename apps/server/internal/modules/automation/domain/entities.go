package domain

type Trigger struct {
	ID         uint
	Name       string
	Event      string
	Conditions []Condition
	Actions    []Action
	Active     bool
}

type Condition struct {
	Field string
	Op    string
	Value interface{}
}

type Action struct {
	Type   string
	Params map[string]interface{}
}

type Execution struct {
	TriggerID uint
	TicketID  uint
	Status    string
	Message   string
}
