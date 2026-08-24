package app

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" format:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type FormStatus string

const (
	FormStatusDraft     FormStatus = "draft"
	FormStatusPublished FormStatus = "published"
)

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeEmail    FieldType = "email"
	FieldTypeNumber   FieldType = "number"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
)

type FormField struct {
	ID          string    `json:"id" minLength:"1"`
	Label       string    `json:"label" minLength:"1"`
	Type        FieldType `json:"type" enum:"text,textarea,email,number,select,checkbox"`
	Required    bool      `json:"required"`
	Placeholder string    `json:"placeholder,omitempty"`
	HelpText    string    `json:"helpText,omitempty"`
	Options     []string  `json:"options,omitempty"`
	Min         *float64  `json:"min,omitempty"`
	Max         *float64  `json:"max,omitempty"`
}

type Form struct {
	ID          string      `json:"id"`
	OwnerID     string      `json:"ownerId"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Slug        string      `json:"slug"`
	Status      FormStatus  `json:"status" enum:"draft,published"`
	Fields      []FormField `json:"fields"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	PublishedAt *time.Time  `json:"publishedAt,omitempty"`
}

type FormSummary struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Slug          string     `json:"slug"`
	Status        FormStatus `json:"status" enum:"draft,published"`
	PublicURL     string     `json:"publicUrl,omitempty"`
	FieldCount    int        `json:"fieldCount"`
	ResponseCount int        `json:"responseCount"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
}

type PublicForm struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Slug        string      `json:"slug"`
	Fields      []FormField `json:"fields"`
}

type FormResponse struct {
	ID          string         `json:"id"`
	FormID      string         `json:"formId"`
	Answers     map[string]any `json:"answers"`
	SubmittedAt time.Time      `json:"submittedAt"`
}
