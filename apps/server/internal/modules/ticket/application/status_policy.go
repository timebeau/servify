package application

import "fmt"

type StatusTransitionPolicy struct{}

func NewStatusTransitionPolicy() StatusTransitionPolicy {
	return StatusTransitionPolicy{}
}

func (p StatusTransitionPolicy) Validate(fromStatus, toStatus string) error {
	if toStatus == "" || fromStatus == toStatus {
		return nil
	}

	switch fromStatus {
	case "", "open":
		if toStatus == "assigned" || toStatus == "resolved" || toStatus == "closed" {
			return nil
		}
	case "assigned":
		if toStatus == "open" || toStatus == "in_progress" || toStatus == "resolved" || toStatus == "closed" {
			return nil
		}
	case "in_progress":
		if toStatus == "open" || toStatus == "resolved" || toStatus == "closed" {
			return nil
		}
	case "resolved", "closed":
		if toStatus == "closed" {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition: %s -> %s", fromStatus, toStatus)
}
