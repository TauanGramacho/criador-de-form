package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Store struct {
	db *Database
}

func NewStore(db *Database) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUserWithPassword(name, email, passwordHash string) (User, error) {
	now := nowUTC()
	user := User{
		ID:        newID("usr"),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Name:      strings.TrimSpace(name),
		CreatedAt: now,
	}
	_, err := s.db.Exec(
		`INSERT INTO users (id, email, name, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID,
		user.Email,
		user.Name,
		passwordHash,
		now.Format(timeFormat),
		now.Format(timeFormat),
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return User{}, ErrConflict
		}
		return User{}, err
	}
	return user, nil
}

func (s *Store) FindUserByEmail(email string) (User, string, error) {
	row := s.db.QueryRow(
		`SELECT id, email, name, COALESCE(password_hash, ''), created_at FROM users WHERE email = ?`,
		strings.ToLower(strings.TrimSpace(email)),
	)
	return scanUserWithPassword(row)
}

func (s *Store) FindUserByID(id string) (User, error) {
	row := s.db.QueryRow(`SELECT id, email, name, created_at FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func (s *Store) UpsertGoogleUser(email, name, googleID string) (User, error) {
	now := nowUTC()
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if name == "" {
		name = email
	}

	existing, _, err := s.FindUserByEmail(email)
	if err == nil {
		_, err = s.db.Exec(
			`UPDATE users SET name = ?, google_id = ?, updated_at = ? WHERE id = ?`,
			name,
			googleID,
			now.Format(timeFormat),
			existing.ID,
		)
		if err != nil {
			return User{}, err
		}
		existing.Name = name
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return User{}, err
	}

	user := User{
		ID:        newID("usr"),
		Email:     email,
		Name:      name,
		CreatedAt: now,
	}
	_, err = s.db.Exec(
		`INSERT INTO users (id, email, name, google_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID,
		user.Email,
		user.Name,
		googleID,
		now.Format(timeFormat),
		now.Format(timeFormat),
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return User{}, ErrConflict
		}
		return User{}, err
	}
	return user, nil
}

func (s *Store) CreateSession(userID, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO sessions (token_hash, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`,
		tokenHash,
		userID,
		expiresAt.UTC().Format(timeFormat),
		nowUTC().Format(timeFormat),
	)
	return err
}

func (s *Store) DeleteSession(tokenHash string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	return err
}

func (s *Store) FindUserBySession(tokenHash string) (User, error) {
	row := s.db.QueryRow(`
		SELECT u.id, u.email, u.name, u.created_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = ? AND s.expires_at > ?
	`, tokenHash, nowUTC().Format(timeFormat))
	return scanUser(row)
}

func (s *Store) CreateForm(ownerID, title, description string, fields []FormField) (Form, error) {
	now := nowUTC()
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return Form{}, err
	}
	form := Form{
		ID:          newID("frm"),
		OwnerID:     ownerID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Slug:        slugify(title),
		Status:      FormStatusDraft,
		Fields:      fields,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err = s.db.Exec(
		`INSERT INTO forms (id, owner_id, title, description, slug, status, fields_json, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		form.ID,
		form.OwnerID,
		form.Title,
		form.Description,
		form.Slug,
		form.Status,
		string(fieldsJSON),
		now.Format(timeFormat),
		now.Format(timeFormat),
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return Form{}, ErrConflict
		}
		return Form{}, err
	}
	return form, nil
}

func (s *Store) ListForms(ownerID, frontendBaseURL string) ([]FormSummary, error) {
	rows, err := s.db.Query(`
		SELECT f.id, f.title, f.description, f.slug, f.status, f.fields_json, f.created_at, f.updated_at, f.published_at,
			COUNT(r.id) AS response_count
		FROM forms f
		LEFT JOIN responses r ON r.form_id = f.id
		WHERE f.owner_id = ?
		GROUP BY f.id
		ORDER BY f.updated_at DESC
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []FormSummary
	for rows.Next() {
		var summary FormSummary
		var fieldsJSON string
		var createdAt, updatedAt string
		var publishedAt sql.NullString
		if err := rows.Scan(
			&summary.ID,
			&summary.Title,
			&summary.Description,
			&summary.Slug,
			&summary.Status,
			&fieldsJSON,
			&createdAt,
			&updatedAt,
			&publishedAt,
			&summary.ResponseCount,
		); err != nil {
			return nil, err
		}
		var fields []FormField
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return nil, err
		}
		summary.FieldCount = len(fields)
		summary.CreatedAt = parseTime(createdAt)
		summary.UpdatedAt = parseTime(updatedAt)
		if publishedAt.Valid {
			t := parseTime(publishedAt.String)
			summary.PublishedAt = &t
		}
		if summary.Status == FormStatusPublished {
			summary.PublicURL = publicFormURL(frontendBaseURL, summary.Slug)
		}
		forms = append(forms, summary)
	}
	return forms, rows.Err()
}

func (s *Store) GetForm(ownerID, id string) (Form, error) {
	row := s.db.QueryRow(`
		SELECT id, owner_id, title, description, slug, status, fields_json, created_at, updated_at, published_at
		FROM forms
		WHERE id = ? AND owner_id = ?
	`, id, ownerID)
	return scanForm(row)
}

func (s *Store) UpdateForm(ownerID, id, title, description string, fields []FormField) (Form, error) {
	now := nowUTC()
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return Form{}, err
	}
	result, err := s.db.Exec(`
		UPDATE forms
		SET title = ?, description = ?, fields_json = ?, updated_at = ?
		WHERE id = ? AND owner_id = ?
	`, strings.TrimSpace(title), strings.TrimSpace(description), string(fieldsJSON), now.Format(timeFormat), id, ownerID)
	if err != nil {
		return Form{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Form{}, err
	}
	if count == 0 {
		return Form{}, ErrNotFound
	}
	return s.GetForm(ownerID, id)
}

func (s *Store) SetFormStatus(ownerID, id string, status FormStatus) (Form, error) {
	now := nowUTC()
	publishedAt := sql.NullString{}
	if status == FormStatusPublished {
		publishedAt = sql.NullString{String: now.Format(timeFormat), Valid: true}
	}
	result, err := s.db.Exec(`
		UPDATE forms
		SET status = ?, published_at = ?, updated_at = ?
		WHERE id = ? AND owner_id = ?
	`, status, publishedAt, now.Format(timeFormat), id, ownerID)
	if err != nil {
		return Form{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Form{}, err
	}
	if count == 0 {
		return Form{}, ErrNotFound
	}
	return s.GetForm(ownerID, id)
}

func (s *Store) GetPublishedFormBySlug(slug string) (PublicForm, error) {
	row := s.db.QueryRow(`
		SELECT id, title, description, slug, fields_json
		FROM forms
		WHERE slug = ? AND status = 'published'
	`, slug)

	var form PublicForm
	var fieldsJSON string
	if err := row.Scan(&form.ID, &form.Title, &form.Description, &form.Slug, &fieldsJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PublicForm{}, ErrNotFound
		}
		return PublicForm{}, err
	}
	if err := json.Unmarshal([]byte(fieldsJSON), &form.Fields); err != nil {
		return PublicForm{}, err
	}
	return form, nil
}

func (s *Store) CreateResponse(formID string, answers map[string]any, submitterIP string) (FormResponse, error) {
	now := nowUTC()
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return FormResponse{}, err
	}
	response := FormResponse{
		ID:          newID("rsp"),
		FormID:      formID,
		Answers:     answers,
		SubmittedAt: now,
	}
	_, err = s.db.Exec(
		`INSERT INTO responses (id, form_id, answers_json, submitted_at, submitter_ip) VALUES (?, ?, ?, ?, ?)`,
		response.ID,
		response.FormID,
		string(answersJSON),
		now.Format(timeFormat),
		submitterIP,
	)
	return response, err
}

func (s *Store) ListResponses(ownerID, formID string) ([]FormResponse, error) {
	rows, err := s.db.Query(`
		SELECT r.id, r.form_id, r.answers_json, r.submitted_at
		FROM responses r
		JOIN forms f ON f.id = r.form_id
		WHERE r.form_id = ? AND f.owner_id = ?
		ORDER BY r.submitted_at DESC
	`, formID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	responses := []FormResponse{}
	for rows.Next() {
		var response FormResponse
		var answersJSON string
		var submittedAt string
		if err := rows.Scan(&response.ID, &response.FormID, &answersJSON, &submittedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(answersJSON), &response.Answers); err != nil {
			return nil, err
		}
		response.SubmittedAt = parseTime(submittedAt)
		responses = append(responses, response)
	}
	return responses, rows.Err()
}

func scanUser(row interface{ Scan(dest ...any) error }) (User, error) {
	var user User
	var createdAt string
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	user.CreatedAt = parseTime(createdAt)
	return user, nil
}

func scanUserWithPassword(row interface{ Scan(dest ...any) error }) (User, string, error) {
	var user User
	var passwordHash string
	var createdAt string
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &passwordHash, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, "", ErrNotFound
		}
		return User{}, "", err
	}
	user.CreatedAt = parseTime(createdAt)
	return user, passwordHash, nil
}

func scanForm(row interface{ Scan(dest ...any) error }) (Form, error) {
	var form Form
	var fieldsJSON string
	var createdAt, updatedAt string
	var publishedAt sql.NullString
	if err := row.Scan(
		&form.ID,
		&form.OwnerID,
		&form.Title,
		&form.Description,
		&form.Slug,
		&form.Status,
		&fieldsJSON,
		&createdAt,
		&updatedAt,
		&publishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Form{}, ErrNotFound
		}
		return Form{}, err
	}
	if err := json.Unmarshal([]byte(fieldsJSON), &form.Fields); err != nil {
		return Form{}, err
	}
	form.CreatedAt = parseTime(createdAt)
	form.UpdatedAt = parseTime(updatedAt)
	if publishedAt.Valid {
		t := parseTime(publishedAt.String)
		form.PublishedAt = &t
	}
	return form, nil
}

func parseTime(value string) time.Time {
	t, err := time.Parse(timeFormat, value)
	if err != nil {
		return time.Time{}
	}
	return t
}

func isUniqueConstraint(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "constraint") || strings.Contains(msg, "unique")
}

func responseCount(db *Database, formID string) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM responses WHERE form_id = ?`, formID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count responses: %w", err)
	}
	return count, nil
}
