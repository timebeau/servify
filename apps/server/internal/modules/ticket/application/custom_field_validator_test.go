package application

import (
	"testing"
)

func TestCustomFieldValidatorValidateRequiredVisibleField(t *testing.T) {
	validator := NewCustomFieldValidator()

	values, err := validator.Validate([]CustomFieldDefinition{
		{
			ID:       1,
			Key:      "severity",
			Type:     "select",
			Required: true,
			Active:   true,
			Options:  []string{"low", "high"},
			ShowWhen: &CustomFieldCondition{
				All: []CustomFieldClause{{Field: "category", Op: "eq", Value: "incident"}},
			},
		},
	}, map[string]interface{}{
		"severity": "high",
	}, map[string]interface{}{
		"category": "incident",
	}, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(values) != 1 || values[0].Value != "high" {
		t.Fatalf("unexpected values: %+v", values)
	}
}

func TestCustomFieldValidatorRejectsMissingRequiredVisibleField(t *testing.T) {
	validator := NewCustomFieldValidator()

	_, err := validator.Validate([]CustomFieldDefinition{
		{
			ID:       1,
			Key:      "severity",
			Type:     "select",
			Required: true,
			Active:   true,
			Options:  []string{"low", "high"},
			ShowWhen: &CustomFieldCondition{
				All: []CustomFieldClause{{Field: "category", Op: "eq", Value: "incident"}},
			},
		},
	}, map[string]interface{}{}, map[string]interface{}{
		"category": "incident",
	}, true)
	if err == nil {
		t.Fatal("expected required field error")
	}
}

func TestCustomFieldValidatorRejectsInvalidOption(t *testing.T) {
	validator := NewCustomFieldValidator()

	_, err := validator.Validate([]CustomFieldDefinition{
		{
			ID:      1,
			Key:     "severity",
			Type:    "select",
			Active:  true,
			Options: []string{"low", "high"},
		},
	}, map[string]interface{}{
		"severity": "critical",
	}, nil, false)
	if err == nil {
		t.Fatal("expected invalid option error")
	}
}

func TestCustomFieldValidatorRejectsNumberOutOfRange(t *testing.T) {
	min := 1.0
	max := 10.0
	validator := NewCustomFieldValidator()

	_, err := validator.Validate([]CustomFieldDefinition{
		{
			ID:     1,
			Key:    "impact_score",
			Type:   "number",
			Active: true,
			Validation: CustomFieldValidation{
				Min: &min,
				Max: &max,
			},
		},
	}, map[string]interface{}{
		"impact_score": 20,
	}, nil, false)
	if err == nil {
		t.Fatal("expected range validation error")
	}
}
