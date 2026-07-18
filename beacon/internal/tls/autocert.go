package tls

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type AutoTLSManager struct {
	manager         *autocert.Manager
	hostname        string
	cacheDir        string
	cancel          context.CancelFunc
	challengeServer *http.Server
}

func NewAutoTLSManager(hostname, cacheDir, email string) *AutoTLSManager {
	var emailContact []string
	if email != "" {
		emailContact = []string{"mailto:" + email}
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(cacheDir),
		HostPolicy: autocert.HostWhitelist(hostname),
		Email:      email,
	}
	_ = emailContact

	return &AutoTLSManager{
		manager:  m,
		hostname: hostname,
		cacheDir: cacheDir,
	}
}

func (m *AutoTLSManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return m.manager.GetCertificate(hello)
}

func (m *AutoTLSManager) HTTPHandler() http.Handler {
	return m.manager.HTTPHandler(nil)
}

func (m *AutoTLSManager) StartChallengeServer() error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: m.HTTPHandler(),
	}
	m.challengeServer = srv
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	go func() {
		_ = srv.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	return nil
}

func (m *AutoTLSManager) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
