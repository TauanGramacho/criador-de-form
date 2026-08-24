package app

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func newID(prefix string) string {
	return prefix + "_" + randomString(18)
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return strings.TrimRight(base64.RawURLEncoding.EncodeToString(b), "=")
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = nonSlugChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "form"
	}
	return fmt.Sprintf("%s-%s", s, strings.ToLower(randomString(5)))
}

func publicFormURL(baseURL, slug string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return "/f/" + url.PathEscape(slug)
	}
	return baseURL + "/f/" + url.PathEscape(slug)
}

func nowUTC() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
