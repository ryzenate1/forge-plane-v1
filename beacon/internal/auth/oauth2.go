package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	OAuth2ExchangeTimeout = 30 * time.Second
	PKCEVerifierLength    = 64
	StateLength           = 32
)

var (
	ErrInvalidState      = errors.New("invalid OAuth2 state parameter")
	ErrInvalidNonce      = errors.New("invalid OAuth2 nonce")
	ErrExchangeTimeout   = errors.New("OAuth2 token exchange timed out")
	ErrInvalidRedirectURL = errors.New("invalid redirect URL")
	ErrTLSRequired       = errors.New("TLS is required for OAuth2 token exchange")
)

type OAuth2Provider struct {
	ClientID        string
	ClientSecret    string
	AuthURL         string
	TokenURL        string
	RedirectURL     string
	Scopes          []string
	allowedDomains  []string
	insecureSkipTLS bool
	config          *oauth2.Config
	httpClient      *http.Client
	mu              sync.Mutex
}

func NewOAuth2Provider(clientID, clientSecret, authURL, tokenURL, redirectURL string, scopes []string, allowedDomains []string) *OAuth2Provider {
	p := &OAuth2Provider{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		AuthURL:        authURL,
		TokenURL:       tokenURL,
		RedirectURL:    redirectURL,
		Scopes:         scopes,
		allowedDomains: allowedDomains,
	}

	p.config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	p.httpClient = &http.Client{
		Timeout: OAuth2ExchangeTimeout,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	return p
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}

func GeneratePKCEVerifier() (string, error) {
	b, err := generateRandomBytes(PKCEVerifierLength)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GeneratePKCEChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func GenerateOAuth2State() (string, error) {
	b, err := generateRandomBytes(StateLength)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (p *OAuth2Provider) WithInsecureSkipTLSVerify() *OAuth2Provider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.insecureSkipTLS = true
	return p
}

func (p *OAuth2Provider) WithAllowedDomains(domains []string) *OAuth2Provider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.allowedDomains = domains
	return p
}

func (p *OAuth2Provider) isRedirectURLAllowed(redirectURL string) bool {
	if len(p.allowedDomains) == 0 {
		return redirectURL == p.RedirectURL
	}
	for _, domain := range p.allowedDomains {
		if strings.Contains(redirectURL, domain) {
			return true
		}
	}
	return false
}

func (p *OAuth2Provider) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	p.mu.Lock()
	config := p.config
	p.mu.Unlock()
	if config == nil {
		config = &oauth2.Config{
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			RedirectURL:  p.RedirectURL,
			Scopes:       p.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  p.AuthURL,
				TokenURL: p.TokenURL,
			},
		}
	}
	return config.AuthCodeURL(state, opts...)
}

func (p *OAuth2Provider) AuthCodeURLWithPKCE(state, codeChallenge string) string {
	return p.AuthCodeURL(state, oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
}

func (p *OAuth2Provider) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if !strings.HasPrefix(p.TokenURL, "https://") && !p.insecureSkipTLS {
		return nil, ErrTLSRequired
	}

	if !p.isRedirectURLAllowed(p.RedirectURL) {
		return nil, ErrInvalidRedirectURL
	}

	ctx, cancel := context.WithTimeout(ctx, OAuth2ExchangeTimeout)
	defer cancel()

	p.mu.Lock()
	config := p.config
	httpClient := p.httpClient
	p.mu.Unlock()

	if config == nil {
		config = &oauth2.Config{
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			RedirectURL:  p.RedirectURL,
			Scopes:       p.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  p.AuthURL,
				TokenURL: p.TokenURL,
			},
		}
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	return config.Exchange(ctx, code, opts...)
}

func (p *OAuth2Provider) ExchangeWithPKCE(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	return p.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
}

func (p *OAuth2Provider) ValidateRedirectURL(redirectURL string) error {
	if !p.isRedirectURLAllowed(redirectURL) {
		return ErrInvalidRedirectURL
	}
	if !strings.HasPrefix(redirectURL, "https://") {
		return ErrTLSRequired
	}
	return nil
}
