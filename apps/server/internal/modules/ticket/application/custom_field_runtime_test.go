package application

import (
	"testing"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/modules/ticket/domain"
)

func TestPrepareCustomFieldMutationTracksDeletesAndUpserts(t *testing.T) {
	fields := []models.CustomField{
		{ID: 1, Key: "severity", Type: "select", Active: true, OptionsJSON: `["low","high"]`},
		{ID: 2, Key: "summary", Type: "string", Active: true},
	}

	mutation, err := PrepareCustomFieldMutation(fields, 42, map[string]interface{}{
		"severity": "high",
		"summary":  nil,
	}, map[string]interface{}{
		"ticket.category": "incident",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mutation == nil {
		t.Fatal("expected mutation")
	}
	if mutation.ClearAll {
		t.Fatal("did not expect clear all")
	}
	if len(mutation.DeleteFieldIDs) != 1 || mutation.DeleteFieldIDs[0] != 2 {
		t.Fatalf("unexpected delete ids: %+v", mutation.DeleteFieldIDs)
	}
	if len(mutation.Upserts) != 1 || mutation.Upserts[0].CustomFieldID != 1 || mutation.Upserts[0].Value != "high" {
		t.Fatalf("unexpected upserts: %+v", mutation.Upserts)
	}
}

func TestPrepareCustomFieldMutationSkipsHiddenFields(t *testing.T) {
	fields := []models.CustomField{
		{
			ID:           1,
			Key:          "severity",
			Type:         "select",
			Active:       true,
			OptionsJSON:  `["low","high"]`,
			ShowWhenJSON: `{"all":[{"field":"ticket.category","op":"eq","value":"incident"}]}`,
		},
	}

	mutation, err := PrepareCustomFieldMutation(fields, 42, map[string]interface{}{
		"severity": "high",
	}, map[string]interface{}{
		"ticket.category": "question",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mutation == nil {
		t.Fatal("expected mutation")
	}
	if len(mutation.Upserts) != 0 {
		t.Fatalf("expected hidden field to be skipped, got %+v", mutation.Upserts)
	}
}

func TestMapMutationToModelValues(t *testing.T) {
	values := MapMutationToModelValues(7, &CustomFieldMutation{
		Upserts: []domain.CustomFieldValue{{
			CustomFieldID: 3,
			Key:           "region",
			Value:         "apac",
		}},
	})
	if len(values) != 1 {
		t.Fatalf("expected one value, got %d", len(values))
	}
	if values[0].TicketID != 7 || values[0].CustomFieldID != 3 || values[0].Value != "apac" {
		t.Fatalf("unexpected values: %+v", values)
	}
}
