package application

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"servify/apps/server/internal/modules/ticket/domain"
)

type CustomFieldDefinition struct {
	ID         uint
	Key        string
	Type       string
	Required   bool
	Active     bool
	Options    []string
	Validation CustomFieldValidation
	ShowWhen   *CustomFieldCondition
}

type CustomFieldValidation struct {
	Min       *float64
	Max       *float64
	MinLength *int
	MaxLength *int
	Regex     string
}

type CustomFieldCondition struct {
	All []CustomFieldClause
	Any []CustomFieldClause
}

type CustomFieldClause struct {
	Field string
	Op    string
	Value interface{}
}

type CustomFieldMutation struct {
	ClearAll       bool
	DeleteFieldIDs []uint
	Upserts        []domain.CustomFieldValue
}

type CustomFieldValidator struct{}

func NewCustomFieldValidator() CustomFieldValidator {
	return CustomFieldValidator{}
}

func (v CustomFieldValidator) Validate(
	definitions []CustomFieldDefinition,
	provided map[string]interface{},
	ticketContext map[string]interface{},
	enforceRequired bool,
) ([]domain.CustomFieldValue, error) {
	fieldByKey := make(map[string]CustomFieldDefinition, len(definitions))
	ctxMap := make(map[string]interface{}, len(ticketContext)+len(provided)+8)
	for k, value := range ticketContext {
		ctxMap[k] = value
	}
	for _, definition := range definitions {
		fieldByKey[definition.Key] = definition
	}
	for k, value := range provided {
		ctxMap[k] = value
		ctxMap["cf."+k] = value
	}

	if enforceRequired {
		for _, definition := range definitions {
			if !definition.Active || !definition.Required {
				continue
			}
			if !customFieldConditionMet(definition.ShowWhen, ctxMap) {
				continue
			}
			raw, ok := provided[definition.Key]
			if !ok || isEmptyCustomFieldValue(raw) {
				return nil, fmt.Errorf("custom field %q is required", definition.Key)
			}
		}
	}

	if len(provided) == 0 {
		return nil, nil
	}

	values := make([]domain.CustomFieldValue, 0, len(provided))
	for key, raw := range provided {
		definition, ok := fieldByKey[key]
		if !ok {
			return nil, fmt.Errorf("unknown custom field: %s", key)
		}
		if !definition.Active || isEmptyCustomFieldValue(raw) || !customFieldConditionMet(definition.ShowWhen, ctxMap) {
			continue
		}

		value, err := normalizeCustomFieldValue(definition, raw)
		if err != nil {
			return nil, fmt.Errorf("custom field %q: %w", key, err)
		}
		values = append(values, domain.CustomFieldValue{
			CustomFieldID: definition.ID,
			Key:           definition.Key,
			Value:         value,
		})
	}

	return values, nil
}

func customFieldConditionMet(condition *CustomFieldCondition, ctx map[string]interface{}) bool {
	if condition == nil {
		return true
	}
	if len(condition.All) > 0 {
		for _, clause := range condition.All {
			if !evalCustomFieldClause(clause, ctx) {
				return false
			}
		}
		return true
	}
	if len(condition.Any) > 0 {
		for _, clause := range condition.Any {
			if evalCustomFieldClause(clause, ctx) {
				return true
			}
		}
		return false
	}
	return true
}

func evalCustomFieldClause(clause CustomFieldClause, ctx map[string]interface{}) bool {
	left, ok := ctx[clause.Field]
	if !ok {
		return false
	}
	switch strings.TrimSpace(clause.Op) {
	case "", "eq":
		return fmt.Sprint(left) == fmt.Sprint(clause.Value)
	case "ne":
		return fmt.Sprint(left) != fmt.Sprint(clause.Value)
	case "in":
		for _, item := range normalizeStringList(clause.Value) {
			if fmt.Sprint(left) == item {
				return true
			}
		}
		return false
	case "not_in":
		for _, item := range normalizeStringList(clause.Value) {
			if fmt.Sprint(left) == item {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func normalizeCustomFieldValue(definition CustomFieldDefinition, raw interface{}) (string, error) {
	switch strings.TrimSpace(definition.Type) {
	case "string":
		value, ok := raw.(string)
		if !ok {
			return "", fmt.Errorf("expected string")
		}
		value = strings.TrimSpace(value)
		if err := validateStringField(definition.Validation, value); err != nil {
			return "", err
		}
		return value, nil
	case "number":
		value, err := normalizeNumber(raw)
		if err != nil {
			return "", err
		}
		if err := validateNumberField(definition.Validation, value); err != nil {
			return "", err
		}
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case "boolean":
		value, err := normalizeBool(raw)
		if err != nil {
			return "", err
		}
		if value {
			return "true", nil
		}
		return "false", nil
	case "date":
		value, ok := raw.(string)
		if !ok {
			return "", fmt.Errorf("expected date string")
		}
		return normalizeDate(value)
	case "select":
		value, ok := raw.(string)
		if !ok {
			return "", fmt.Errorf("expected string")
		}
		value = strings.TrimSpace(value)
		if err := validateOption(definition.Options, value, false); err != nil {
			return "", err
		}
		return value, nil
	case "multiselect":
		values, err := normalizeStringListAny(raw)
		if err != nil {
			return "", err
		}
		for _, value := range values {
			if err := validateOption(definition.Options, value, false); err != nil {
				return "", err
			}
		}
		return strings.Join(values, ","), nil
	default:
		return "", fmt.Errorf("unsupported type: %s", definition.Type)
	}
}

func validateStringField(rule CustomFieldValidation, value string) error {
	if rule.MinLength != nil && len(value) < *rule.MinLength {
		return fmt.Errorf("length must be >= %d", *rule.MinLength)
	}
	if rule.MaxLength != nil && len(value) > *rule.MaxLength {
		return fmt.Errorf("length must be <= %d", *rule.MaxLength)
	}
	if strings.TrimSpace(rule.Regex) != "" {
		pattern, err := regexp.Compile(rule.Regex)
		if err != nil {
			return fmt.Errorf("invalid regex")
		}
		if !pattern.MatchString(value) {
			return fmt.Errorf("does not match pattern")
		}
	}
	return nil
}

func validateNumberField(rule CustomFieldValidation, value float64) error {
	if rule.Min != nil && value < *rule.Min {
		return fmt.Errorf("must be >= %v", *rule.Min)
	}
	if rule.Max != nil && value > *rule.Max {
		return fmt.Errorf("must be <= %v", *rule.Max)
	}
	return nil
}

func validateOption(options []string, value string, allowEmpty bool) error {
	value = strings.TrimSpace(value)
	if value == "" {
		if allowEmpty {
			return nil
		}
		return fmt.Errorf("value required")
	}
	if len(options) == 0 || slices.Contains(options, value) {
		return nil
	}
	return fmt.Errorf("invalid option")
}

func normalizeNumber(value interface{}) (float64, error) {
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case float32:
		return float64(typed), nil
	case int:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case int32:
		return float64(typed), nil
	case uint:
		return float64(typed), nil
	case uint64:
		return float64(typed), nil
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return 0, fmt.Errorf("empty number")
		}
		number, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number")
		}
		return number, nil
	default:
		return 0, fmt.Errorf("invalid number")
	}
}

func normalizeBool(value interface{}) (bool, error) {
	switch typed := value.(type) {
	case bool:
		return typed, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "1", "yes", "y":
			return true, nil
		case "false", "0", "no", "n":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean")
		}
	default:
		return false, fmt.Errorf("invalid boolean")
	}
}

func normalizeDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if len(value) == 10 {
		if _, err := time.Parse("2006-01-02", value); err == nil {
			return value, nil
		}
	}
	date, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", fmt.Errorf("invalid date")
	}
	return date.Format("2006-01-02"), nil
}

func normalizeStringListAny(value interface{}) ([]string, error) {
	switch typed := value.(type) {
	case []string:
		return normalizeStringItems(typed), nil
	case []interface{}:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, fmt.Sprint(item))
		}
		return normalizeStringItems(items), nil
	case string:
		return normalizeStringItems(strings.Split(typed, ",")), nil
	default:
		return nil, fmt.Errorf("expected string list")
	}
}

func normalizeStringList(value interface{}) []string {
	values, err := normalizeStringListAny(value)
	if err != nil {
		return nil
	}
	return values
}

func normalizeStringItems(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func isEmptyCustomFieldValue(value interface{}) bool {
	if value == nil {
		return true
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed) == ""
	case []string:
		return len(normalizeStringItems(typed)) == 0
	case []interface{}:
		values, err := normalizeStringListAny(typed)
		return err != nil || len(values) == 0
	default:
		return strings.TrimSpace(fmt.Sprint(value)) == ""
	}
}
