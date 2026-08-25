package app

import (
	"fmt"
	"net/mail"
	"strings"
)

func validateFormFields(fields []FormField) error {
	if len(fields) == 0 {
		return fmt.Errorf("form must contain at least one field")
	}
	seen := map[string]bool{}
	for i, field := range fields {
		if strings.TrimSpace(field.ID) == "" {
			return fmt.Errorf("field %d must have an id", i+1)
		}
		if seen[field.ID] {
			return fmt.Errorf("field id %q is duplicated", field.ID)
		}
		seen[field.ID] = true
		if strings.TrimSpace(field.Label) == "" {
			return fmt.Errorf("field %q must have a label", field.ID)
		}
		switch field.Type {
		case FieldTypeText, FieldTypeTextarea, FieldTypeEmail, FieldTypeNumber, FieldTypeCheckbox:
		case FieldTypeSelect:
			if len(field.Options) == 0 {
				return fmt.Errorf("select field %q must have options", field.ID)
			}
			optionSeen := map[string]bool{}
			for _, option := range field.Options {
				option = strings.TrimSpace(option)
				if option == "" {
					return fmt.Errorf("select field %q has an empty option", field.ID)
				}
				if optionSeen[option] {
					return fmt.Errorf("select field %q has duplicated option %q", field.ID, option)
				}
				optionSeen[option] = true
			}
		default:
			return fmt.Errorf("field %q has unsupported type %q", field.ID, field.Type)
		}
	}
	return nil
}

func validateAnswers(fields []FormField, answers map[string]any) error {
	fieldByID := make(map[string]FormField, len(fields))
	for _, field := range fields {
		fieldByID[field.ID] = field
	}
	for key := range answers {
		if _, ok := fieldByID[key]; !ok {
			return fmt.Errorf("answer for unknown field %q", key)
		}
	}
	for _, field := range fields {
		value, exists := answers[field.ID]
		if !exists || isEmptyAnswer(value) {
			if field.Required {
				return fmt.Errorf("field %q is required", field.ID)
			}
			continue
		}
		if err := validateAnswer(field, value); err != nil {
			return err
		}
	}
	return nil
}

func validateAnswer(field FormField, value any) error {
	switch field.Type {
	case FieldTypeText, FieldTypeTextarea:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field %q must be text", field.ID)
		}
	case FieldTypeEmail:
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("field %q must be an email", field.ID)
		}
		if _, err := mail.ParseAddress(text); err != nil {
			return fmt.Errorf("field %q must contain a valid email", field.ID)
		}
	case FieldTypeNumber:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("field %q must be a number", field.ID)
		}
	case FieldTypeSelect:
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("field %q must be one of the configured options", field.ID)
		}
		for _, option := range field.Options {
			if text == option {
				return nil
			}
		}
		return fmt.Errorf("field %q must be one of the configured options", field.ID)
	case FieldTypeCheckbox:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field %q must be true or false", field.ID)
		}
	default:
		return fmt.Errorf("field %q has unsupported type %q", field.ID, field.Type)
	}
	return nil
}

func isEmptyAnswer(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	default:
		return false
	}
}
