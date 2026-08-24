package app

import "testing"

func TestValidateAnswersRejectsInvalidEmailAndUnknownFields(t *testing.T) {
	fields := []FormField{
		{ID: "email", Label: "E-mail", Type: FieldTypeEmail, Required: true},
	}

	err := validateAnswers(fields, map[string]any{
		"email": "not-an-email",
	})
	if err == nil {
		t.Fatal("expected invalid email error")
	}

	err = validateAnswers(fields, map[string]any{
		"email": "user@example.com",
		"extra": "value",
	})
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestValidateAnswersAcceptsConfiguredTypes(t *testing.T) {
	minValue := 1.0
	maxValue := 10.0
	fields := []FormField{
		{ID: "name", Label: "Nome", Type: FieldTypeText, Required: true},
		{ID: "email", Label: "E-mail", Type: FieldTypeEmail, Required: true},
		{ID: "score", Label: "Nota", Type: FieldTypeNumber, Required: true, Min: &minValue, Max: &maxValue},
		{ID: "plan", Label: "Plano", Type: FieldTypeSelect, Options: []string{"Free", "Pro"}},
		{ID: "accepted", Label: "Aceite", Type: FieldTypeCheckbox, Required: true},
	}

	err := validateAnswers(fields, map[string]any{
		"name":     "Maria",
		"email":    "maria@example.com",
		"score":    9.0,
		"plan":     "Pro",
		"accepted": true,
	})
	if err != nil {
		t.Fatalf("expected valid answers, got %v", err)
	}
}

func TestValidateFormFieldsRequiresSelectOptions(t *testing.T) {
	err := validateFormFields([]FormField{
		{ID: "plan", Label: "Plano", Type: FieldTypeSelect, Required: false},
	})
	if err == nil {
		t.Fatal("expected select options error")
	}
}
