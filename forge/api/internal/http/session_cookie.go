package http

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SessionCookieName = "__Host-forge_session"
	CSRFCookieName    = "__Host-forge_csrf"
)

type SessionCookieConfig struct {
	Secure   bool
	SameSite http.SameSite
}

func LoadSessionCookieConfig() SessionCookieConfig {
	secure := true
	if v := strings.ToLower(os.Getenv("SESSION_COOKIE_SECURE")); v == "false" || v == "0" {
		secure = false
	}
	sameSite := http.SameSiteLaxMode
	if v := strings.ToLower(os.Getenv("SESSION_COOKIE_SAME_SITE")); v == "strict" {
		sameSite = http.SameSiteStrictMode
	} else if v == "none" {
		sameSite = http.SameSiteNoneMode
	}
	return SessionCookieConfig{Secure: secure, SameSite: sameSite}
}

func ValidateSessionCookieConfig(cfg SessionCookieConfig, appEnv string) error {
	if appEnv == "production" && !cfg.Secure {
		return ErrInsecureCookieConfig
	}
	return nil
}

func generateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

var ErrInsecureCookieConfig = errors.New("SESSION_COOKIE_SECURE must be true in production")

func setSessionCookie(w http.ResponseWriter, token string, expires time.Time, cfg SessionCookieConfig) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
	}
	http.SetCookie(w, cookie)
}

func setCSRFCookie(w http.ResponseWriter, token string, expires time.Time, cfg SessionCookieConfig) {
	cookie := &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
	}
	http.SetCookie(w, cookie)
}

func clearSessionCookie(w http.ResponseWriter, cfg SessionCookieConfig) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}

func clearCSRFCookie(w http.ResponseWriter, cfg SessionCookieConfig) {
	cookie := &http.Cookie{
		Name:     CSRFCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}

func getCSRFTokenFromCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func getSessionTokenFromCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}
