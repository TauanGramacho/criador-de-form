package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const sessionCookieName = "fb_session"

type Service struct {
	cfg   Config
	store *Store
	oauth *oauth2.Config
}

func NewService(cfg Config, store *Store) *Service {
	var oauthConfig *oauth2.Config
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		oauthConfig = &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}
	return &Service{cfg: cfg, store: store, oauth: oauthConfig}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func verifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Service) createSession(userID string) (http.Cookie, error) {
	token := randomString(36)
	expiresAt := nowUTC().Add(14 * 24 * time.Hour)
	if err := s.store.CreateSession(userID, hashToken(token), expiresAt); err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.CookieSecure,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	}, nil
}

func (s *Service) requireUser(sessionToken string) (User, error) {
	if sessionToken == "" {
		return User{}, huma.Error401Unauthorized("authentication required")
	}
	user, err := s.store.FindUserBySession(hashToken(sessionToken))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, huma.Error401Unauthorized("authentication required")
		}
		return User{}, err
	}
	return user, nil
}

func (s *Service) clearSessionCookie() http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.CookieSecure,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
}

func (s *Service) exchangeGoogleCode(ctx context.Context, code string) (googleUserInfo, error) {
	if s.oauth == nil {
		return googleUserInfo{}, huma.Error400BadRequest("Google authentication is not configured")
	}
	token, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return googleUserInfo{}, huma.Error400BadRequest("invalid Google OAuth code")
	}

	client := s.oauth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return googleUserInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return googleUserInfo{}, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return googleUserInfo{}, err
	}
	if strings.TrimSpace(info.Email) == "" {
		return googleUserInfo{}, huma.Error400BadRequest("Google account did not return an email")
	}
	return info, nil
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
}

func frontendRedirect(baseURL, redirectPath string) string {
	if !strings.HasPrefix(redirectPath, "/") || strings.HasPrefix(redirectPath, "//") {
		redirectPath = "/"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return redirectPath
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL + redirectPath
	}
	rel, err := url.Parse(redirectPath)
	if err != nil {
		return baseURL + "/"
	}
	return u.ResolveReference(rel).String()
}
