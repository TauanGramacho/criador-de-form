package app

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type SessionInput struct {
	Session string `cookie:"fb_session"`
}

type RegisterRequest struct {
	Name            string `json:"name" minLength:"1" maxLength:"120"`
	Email           string `json:"email" format:"email"`
	Password        string `json:"password" minLength:"8" maxLength:"120"`
	ConfirmPassword string `json:"confirmPassword" minLength:"8" maxLength:"120"`
}

type LoginRequest struct {
	Email    string `json:"email" format:"email"`
	Password string `json:"password" minLength:"1"`
}

type AuthResponse struct {
	User User `json:"user"`
}

type AuthConfigResponse struct {
	GoogleEnabled bool `json:"googleEnabled"`
}

type AuthOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      AuthResponse
}

type RegisterOutput struct {
	Body AuthResponse
}

type MeOutput struct {
	Body AuthResponse
}

type AuthConfigOutput struct {
	Body AuthConfigResponse
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type CreateFormRequest struct {
	Title       string      `json:"title" minLength:"1" maxLength:"160"`
	Description string      `json:"description,omitempty" maxLength:"2000"`
	Fields      []FormField `json:"fields"`
}

type UpdateFormRequest = CreateFormRequest

type FormDetail struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Slug          string      `json:"slug"`
	Status        FormStatus  `json:"status" enum:"draft,published"`
	PublicURL     string      `json:"publicUrl,omitempty"`
	Fields        []FormField `json:"fields"`
	ResponseCount int         `json:"responseCount"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
	PublishedAt   string      `json:"publishedAt,omitempty"`
}

type FormDetailOutput struct {
	Body FormDetail
}

type ListFormsOutput struct {
	Body []FormSummary
}

type ListResponsesOutput struct {
	Body []FormResponse
}

type PublicFormOutput struct {
	Body PublicForm
}

type SubmitResponseRequest struct {
	Answers map[string]any `json:"answers"`
}

type SubmitResponseResult struct {
	ID          string `json:"id"`
	SubmittedAt string `json:"submittedAt"`
}

type SubmitResponseOutput struct {
	Body SubmitResponseResult
}

type GoogleStartOutput struct {
	Location  string        `header:"Location"`
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

type GoogleCallbackOutput struct {
	Location  string        `header:"Location"`
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func NewHandler(cfg Config, store *Store) (http.Handler, huma.API) {
	router := chi.NewRouter()
	config := huma.DefaultConfig("Form Builder API", "1.0.0")
	config.Info.Description = "API for the Full Stack Developer technical challenge form builder."
	api := humachi.New(router, config)

	service := NewService(cfg, store)
	service.RegisterRoutes(api)
	return cors(cfg, frontendFallback(router, cfg.FrontendDist)), api
}

func (s *Service) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "registerWithEmail",
		Method:        http.MethodPost,
		Path:          "/api/auth/register",
		Summary:       "Create an account with e-mail and password",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusCreated,
	}, s.registerWithEmail)

	huma.Register(api, huma.Operation{
		OperationID: "loginWithEmail",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Summary:     "Log in with e-mail and password",
		Tags:        []string{"Auth"},
	}, s.loginWithEmail)

	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Summary:     "Log out the current user",
		Tags:        []string{"Auth"},
	}, s.logout)

	huma.Register(api, huma.Operation{
		OperationID: "getCurrentUser",
		Method:      http.MethodGet,
		Path:        "/api/auth/me",
		Summary:     "Get the authenticated user",
		Tags:        []string{"Auth"},
	}, s.me)

	huma.Register(api, huma.Operation{
		OperationID: "getAuthConfig",
		Method:      http.MethodGet,
		Path:        "/api/auth/config",
		Summary:     "Get public authentication provider configuration",
		Tags:        []string{"Auth"},
	}, s.authConfig)

	huma.Register(api, huma.Operation{
		OperationID:   "startGoogleLogin",
		Method:        http.MethodGet,
		Path:          "/api/auth/google/start",
		Summary:       "Start Google OAuth login",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusFound,
	}, s.startGoogleLogin)

	huma.Register(api, huma.Operation{
		OperationID:   "completeGoogleLogin",
		Method:        http.MethodGet,
		Path:          "/api/auth/google/callback",
		Summary:       "Complete Google OAuth login",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusFound,
	}, s.completeGoogleLogin)

	huma.Register(api, huma.Operation{
		OperationID: "listForms",
		Method:      http.MethodGet,
		Path:        "/api/forms",
		Summary:     "List forms for the authenticated administrator",
		Tags:        []string{"Forms"},
	}, s.listForms)

	huma.Register(api, huma.Operation{
		OperationID:   "createForm",
		Method:        http.MethodPost,
		Path:          "/api/forms",
		Summary:       "Create a form",
		Tags:          []string{"Forms"},
		DefaultStatus: http.StatusCreated,
	}, s.createForm)

	huma.Register(api, huma.Operation{
		OperationID: "getForm",
		Method:      http.MethodGet,
		Path:        "/api/forms/{id}",
		Summary:     "Get a form",
		Tags:        []string{"Forms"},
	}, s.getForm)

	huma.Register(api, huma.Operation{
		OperationID: "updateForm",
		Method:      http.MethodPut,
		Path:        "/api/forms/{id}",
		Summary:     "Update a form",
		Tags:        []string{"Forms"},
	}, s.updateForm)

	huma.Register(api, huma.Operation{
		OperationID: "publishForm",
		Method:      http.MethodPost,
		Path:        "/api/forms/{id}/publish",
		Summary:     "Publish a form",
		Tags:        []string{"Forms"},
	}, s.publishForm)

	huma.Register(api, huma.Operation{
		OperationID: "unpublishForm",
		Method:      http.MethodPost,
		Path:        "/api/forms/{id}/unpublish",
		Summary:     "Unpublish a form",
		Tags:        []string{"Forms"},
	}, s.unpublishForm)

	huma.Register(api, huma.Operation{
		OperationID: "listFormResponses",
		Method:      http.MethodGet,
		Path:        "/api/forms/{id}/responses",
		Summary:     "List responses received by a form",
		Tags:        []string{"Responses"},
	}, s.listResponses)

	huma.Register(api, huma.Operation{
		OperationID: "getPublicForm",
		Method:      http.MethodGet,
		Path:        "/api/public/forms/{slug}",
		Summary:     "Get a published public form",
		Tags:        []string{"Public"},
	}, s.getPublicForm)

	huma.Register(api, huma.Operation{
		OperationID:   "submitPublicFormResponse",
		Method:        http.MethodPost,
		Path:          "/api/public/forms/{slug}/responses",
		Summary:       "Submit a response to a published form",
		Tags:          []string{"Public"},
		DefaultStatus: http.StatusCreated,
	}, s.submitPublicResponse)
}

func (s *Service) registerWithEmail(ctx context.Context, input *struct{ Body RegisterRequest }) (*RegisterOutput, error) {
	if input.Body.Password != input.Body.ConfirmPassword {
		return nil, huma.Error400BadRequest("password confirmation does not match")
	}
	passwordHash, err := hashPassword(input.Body.Password)
	if err != nil {
		return nil, err
	}
	user, err := s.store.CreateUserWithPassword(input.Body.Name, input.Body.Email, passwordHash)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return nil, huma.Error409Conflict("email is already registered")
		}
		return nil, err
	}
	return &RegisterOutput{Body: AuthResponse{User: user}}, nil
}

func (s *Service) loginWithEmail(ctx context.Context, input *struct{ Body LoginRequest }) (*AuthOutput, error) {
	user, passwordHash, err := s.store.FindUserByEmail(input.Body.Email)
	if err != nil || passwordHash == "" || !verifyPassword(passwordHash, input.Body.Password) {
		return nil, huma.Error401Unauthorized("invalid email or password")
	}
	cookie, err := s.createSession(user.ID)
	if err != nil {
		return nil, err
	}
	return &AuthOutput{SetCookie: cookie, Body: AuthResponse{User: user}}, nil
}

func (s *Service) logout(ctx context.Context, input *SessionInput) (*LogoutOutput, error) {
	if input.Session != "" {
		_ = s.store.DeleteSession(hashToken(input.Session))
	}
	return &LogoutOutput{SetCookie: s.clearSessionCookie()}, nil
}

func (s *Service) me(ctx context.Context, input *SessionInput) (*MeOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	return &MeOutput{Body: AuthResponse{User: user}}, nil
}

func (s *Service) authConfig(ctx context.Context, input *struct{}) (*AuthConfigOutput, error) {
	return &AuthConfigOutput{
		Body: AuthConfigResponse{
			GoogleEnabled: s.oauth != nil,
		},
	}, nil
}

func (s *Service) startGoogleLogin(ctx context.Context, input *struct {
	Redirect string `query:"redirect"`
}) (*GoogleStartOutput, error) {
	if s.oauth == nil {
		return nil, huma.Error400BadRequest("Google authentication is not configured")
	}
	state := randomString(32)
	redirectPath := input.Redirect
	if redirectPath == "" {
		redirectPath = "/"
	}
	location := s.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return &GoogleStartOutput{
		Location: location,
		SetCookie: []http.Cookie{
			shortCookie("fb_google_state", state, s.cfg.CookieSecure),
			shortCookie("fb_google_redirect", redirectPath, s.cfg.CookieSecure),
		},
	}, nil
}

func (s *Service) completeGoogleLogin(ctx context.Context, input *struct {
	Code           string `query:"code"`
	State          string `query:"state"`
	StateCookie    string `cookie:"fb_google_state"`
	RedirectCookie string `cookie:"fb_google_redirect"`
}) (*GoogleCallbackOutput, error) {
	if input.State == "" || input.StateCookie == "" || input.State != input.StateCookie {
		return nil, huma.Error400BadRequest("invalid Google OAuth state")
	}
	info, err := s.exchangeGoogleCode(ctx, input.Code)
	if err != nil {
		return nil, err
	}
	user, err := s.store.UpsertGoogleUser(info.Email, info.Name, info.ID)
	if err != nil {
		return nil, err
	}
	sessionCookie, err := s.createSession(user.ID)
	if err != nil {
		return nil, err
	}
	return &GoogleCallbackOutput{
		Location: frontendRedirect(s.cfg.FrontendBaseURL, input.RedirectCookie),
		SetCookie: []http.Cookie{
			sessionCookie,
			clearCookie("fb_google_state", s.cfg.CookieSecure),
			clearCookie("fb_google_redirect", s.cfg.CookieSecure),
		},
	}, nil
}

func (s *Service) listForms(ctx context.Context, input *SessionInput) (*ListFormsOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	forms, err := s.store.ListForms(user.ID, s.cfg.FrontendBaseURL)
	if err != nil {
		return nil, err
	}
	return &ListFormsOutput{Body: forms}, nil
}

func (s *Service) createForm(ctx context.Context, input *struct {
	SessionInput
	Body CreateFormRequest
}) (*FormDetailOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	if err := validateFormFields(input.Body.Fields); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	form, err := s.store.CreateForm(user.ID, input.Body.Title, input.Body.Description, input.Body.Fields)
	if err != nil {
		return nil, err
	}
	return s.formDetailOutput(form)
}

func (s *Service) getForm(ctx context.Context, input *struct {
	SessionInput
	ID string `path:"id"`
}) (*FormDetailOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	form, err := s.store.GetForm(user.ID, input.ID)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return s.formDetailOutput(form)
}

func (s *Service) updateForm(ctx context.Context, input *struct {
	SessionInput
	ID   string `path:"id"`
	Body UpdateFormRequest
}) (*FormDetailOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	if err := validateFormFields(input.Body.Fields); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	form, err := s.store.UpdateForm(user.ID, input.ID, input.Body.Title, input.Body.Description, input.Body.Fields)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return s.formDetailOutput(form)
}

func (s *Service) publishForm(ctx context.Context, input *struct {
	SessionInput
	ID string `path:"id"`
}) (*FormDetailOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	form, err := s.store.SetFormStatus(user.ID, input.ID, FormStatusPublished)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return s.formDetailOutput(form)
}

func (s *Service) unpublishForm(ctx context.Context, input *struct {
	SessionInput
	ID string `path:"id"`
}) (*FormDetailOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	form, err := s.store.SetFormStatus(user.ID, input.ID, FormStatusDraft)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return s.formDetailOutput(form)
}

func (s *Service) listResponses(ctx context.Context, input *struct {
	SessionInput
	ID string `path:"id"`
}) (*ListResponsesOutput, error) {
	user, err := s.requireUser(input.Session)
	if err != nil {
		return nil, err
	}
	responses, err := s.store.ListResponses(user.ID, input.ID)
	if err != nil {
		return nil, err
	}
	return &ListResponsesOutput{Body: responses}, nil
}

func (s *Service) getPublicForm(ctx context.Context, input *struct {
	Slug string `path:"slug"`
}) (*PublicFormOutput, error) {
	form, err := s.store.GetPublishedFormBySlug(input.Slug)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return &PublicFormOutput{Body: form}, nil
}

func (s *Service) submitPublicResponse(ctx context.Context, input *struct {
	Slug          string `path:"slug"`
	XForwardedFor string `header:"X-Forwarded-For"`
	Body          SubmitResponseRequest
}) (*SubmitResponseOutput, error) {
	form, err := s.store.GetPublishedFormBySlug(input.Slug)
	if err != nil {
		return nil, mapStoreError(err)
	}
	if err := validateAnswers(form.Fields, input.Body.Answers); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	submitterIP := strings.Split(input.XForwardedFor, ",")[0]
	response, err := s.store.CreateResponse(form.ID, input.Body.Answers, strings.TrimSpace(submitterIP))
	if err != nil {
		return nil, err
	}
	return &SubmitResponseOutput{
		Body: SubmitResponseResult{
			ID:          response.ID,
			SubmittedAt: response.SubmittedAt.Format(timeFormat),
		},
	}, nil
}

func (s *Service) formDetailOutput(form Form) (*FormDetailOutput, error) {
	count, err := responseCount(s.store.db, form.ID)
	if err != nil {
		return nil, err
	}
	detail := FormDetail{
		ID:            form.ID,
		Title:         form.Title,
		Description:   form.Description,
		Slug:          form.Slug,
		Status:        form.Status,
		Fields:        form.Fields,
		ResponseCount: count,
		CreatedAt:     form.CreatedAt.Format(timeFormat),
		UpdatedAt:     form.UpdatedAt.Format(timeFormat),
	}
	if form.Status == FormStatusPublished {
		detail.PublicURL = publicFormURL(s.cfg.FrontendBaseURL, form.Slug)
	}
	if form.PublishedAt != nil {
		detail.PublishedAt = form.PublishedAt.Format(timeFormat)
	}
	return &FormDetailOutput{Body: detail}, nil
}

func mapStoreError(err error) error {
	if errors.Is(err, ErrNotFound) {
		return huma.Error404NotFound("resource not found")
	}
	if errors.Is(err, ErrConflict) {
		return huma.Error409Conflict("resource already exists")
	}
	return err
}

func shortCookie(name, value string, secure bool) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   10 * 60,
		Expires:  nowUTC().Add(10 * time.Minute),
	}
}

func clearCookie(name string, secure bool) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
}

func cors(cfg Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (origin == cfg.FrontendOrigin || cfg.FrontendOrigin == "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
