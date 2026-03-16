package application

import (
	"context"
	"fmt"
)

type CustomerLookupPort interface {
	GetCustomerSummary(ctx context.Context, customerID uint) (map[string]interface{}, error)
}

type CustomerLookupTool struct {
	port CustomerLookupPort
}

func NewCustomerLookupTool(port CustomerLookupPort) *CustomerLookupTool {
	return &CustomerLookupTool{port: port}
}

func (t *CustomerLookupTool) Name() string { return "customer_lookup" }

func (t *CustomerLookupTool) Description() string {
	return "Lookup a customer summary by customer id."
}

func (t *CustomerLookupTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"customer_id": map[string]interface{}{
				"type": "integer",
			},
		},
		"required": []string{"customer_id"},
	}
}

func (t *CustomerLookupTool) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if t.port == nil {
		return nil, fmt.Errorf("customer lookup port is not configured")
	}
	id, err := getUintInput(input, "customer_id")
	if err != nil {
		return nil, err
	}
	return t.port.GetCustomerSummary(ctx, id)
}
