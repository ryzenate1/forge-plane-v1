package models

import (
	"fmt"
	"net"
	"time"
)

type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusCrashed  ServerStatus = "crashed"
)

func ValidServerStatus(s ServerStatus) bool {
	switch s {
	case ServerStatusStarting, ServerStatusRunning, ServerStatusStopping, ServerStatusStopped, ServerStatusCrashed:
		return true
	}
	return false
}

type NodeStatus string

const (
	NodeStatusOnline       NodeStatus = "online"
	NodeStatusOffline      NodeStatus = "offline"
	NodeStatusMaintenance  NodeStatus = "maintenance"
)

func ValidNodeStatus(s NodeStatus) bool {
	switch s {
	case NodeStatusOnline, NodeStatusOffline, NodeStatusMaintenance:
		return true
	}
	return false
}

type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusRunning   BackupStatus = "running"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
)

func ValidBackupStatus(s BackupStatus) bool {
	switch s {
	case BackupStatusPending, BackupStatusRunning, BackupStatusCompleted, BackupStatusFailed:
		return true
	}
	return false
}

type Server struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	NodeID    string       `json:"nodeId"`
	Status    ServerStatus `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (s *Server) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("server id is required")
	}
	if s.Name == "" {
		return fmt.Errorf("server name is required")
	}
	if s.NodeID == "" {
		return fmt.Errorf("server node_id is required")
	}
	if !ValidServerStatus(s.Status) {
		return fmt.Errorf("invalid server status: %q", s.Status)
	}
	if s.CreatedAt.IsZero() {
		return fmt.Errorf("server created_at is required")
	}
	if s.UpdatedAt.IsZero() {
		return fmt.Errorf("server updated_at is required")
	}
	return nil
}

type Node struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	IPAddress string       `json:"ipAddress"`
	Status    NodeStatus   `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (n *Node) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("node id is required")
	}
	if n.Name == "" {
		return fmt.Errorf("node name is required")
	}
	if n.IPAddress == "" {
		return fmt.Errorf("node ip_address is required")
	}
	if net.ParseIP(n.IPAddress) == nil {
		return fmt.Errorf("invalid node ip_address: %q", n.IPAddress)
	}
	if !ValidNodeStatus(n.Status) {
		return fmt.Errorf("invalid node status: %q", n.Status)
	}
	if n.CreatedAt.IsZero() {
		return fmt.Errorf("node created_at is required")
	}
	if n.UpdatedAt.IsZero() {
		return fmt.Errorf("node updated_at is required")
	}
	return nil
}

type Backup struct {
	ID           string       `json:"id"`
	ServerID     string       `json:"serverId"`
	Path         string       `json:"path"`
	SizeBytes    int64        `json:"sizeBytes"`
	Checksum     string       `json:"checksum"`
	Status       BackupStatus `json:"status"`
	ErrorMessage string       `json:"errorMessage,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

func (b *Backup) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("backup id is required")
	}
	if b.ServerID == "" {
		return fmt.Errorf("backup server_id is required")
	}
	if b.Path == "" {
		return fmt.Errorf("backup path is required")
	}
	if !ValidBackupStatus(b.Status) {
		return fmt.Errorf("invalid backup status: %q", b.Status)
	}
	if b.CreatedAt.IsZero() {
		return fmt.Errorf("backup created_at is required")
	}
	if b.UpdatedAt.IsZero() {
		return fmt.Errorf("backup updated_at is required")
	}
	return nil
}
