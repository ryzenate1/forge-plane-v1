package store

import (
	"context"
	"errors"
)

// GetNodeDaemonToken returns the node's current daemon bearer token (the part
// after the `.` separator in `<id>.<token>`). It is reserved for onboarding
// responses that must configure the daemon.
func (s *Store) GetNodeDaemonToken(ctx context.Context, nodeID string) (string, error) {
	if s.db == nil {
		return "", errors.New("no database connection")
	}
	var plaintext, encrypted string
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(daemon_token, ''), COALESCE(daemon_token_encrypted, '')
		FROM nodes
		WHERE id::text = $1
	`, nodeID).Scan(&plaintext, &encrypted)
	if err != nil {
		return "", err
	}
	return s.decryptSecret(encrypted, plaintext, secretAAD("nodes", nodeID, "daemon_token"))
}

// GetNodeDaemonCredential returns the complete credential used to sign panel
// requests to a node. Unlike GetNodeDaemonToken, this is not an onboarding DTO.
func (s *Store) GetNodeDaemonCredential(ctx context.Context, nodeID string) (string, error) {
	if s.db == nil {
		return "", errors.New("no database connection")
	}
	var tokenID string
	if err := s.db.QueryRow(ctx, `SELECT COALESCE(daemon_token_id, '') FROM nodes WHERE id::text = $1`, nodeID).Scan(&tokenID); err != nil {
		return "", err
	}
	token, err := s.GetNodeDaemonToken(ctx, nodeID)
	if err != nil || tokenID == "" || token == "" {
		return "", err
	}
	return tokenID + "." + token, nil
}
