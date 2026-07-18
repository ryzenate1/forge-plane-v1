package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerModel(t *testing.T) {
	now := time.Now()

	server := Server{
		ID:        "1",
		Name:      "Test Server",
		NodeID:    "node1",
		Status:    ServerStatusRunning,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "1", server.ID)
	assert.Equal(t, "Test Server", server.Name)
	assert.Equal(t, "node1", server.NodeID)
	assert.Equal(t, ServerStatusRunning, server.Status)
	assert.Equal(t, now, server.CreatedAt)
	assert.Equal(t, now, server.UpdatedAt)
}

func TestServerStatus(t *testing.T) {
	testCases := []struct {
		status      ServerStatus
		stringValue string
	}{
		{ServerStatusStarting, "starting"},
		{ServerStatusRunning, "running"},
		{ServerStatusStopping, "stopping"},
		{ServerStatusStopped, "stopped"},
		{ServerStatusCrashed, "crashed"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.stringValue, string(tc.status))
	}
}

func TestValidServerStatus(t *testing.T) {
	assert.True(t, ValidServerStatus(ServerStatusStarting))
	assert.True(t, ValidServerStatus(ServerStatusRunning))
	assert.True(t, ValidServerStatus(ServerStatusStopping))
	assert.True(t, ValidServerStatus(ServerStatusStopped))
	assert.True(t, ValidServerStatus(ServerStatusCrashed))
	assert.False(t, ValidServerStatus(ServerStatus("invalid")))
	assert.False(t, ValidServerStatus(ServerStatus("")))
}

func TestServerValidate(t *testing.T) {
	now := time.Now()

	t.Run("valid server", func(t *testing.T) {
		s := Server{
			ID:        "1",
			Name:      "Test",
			NodeID:    "node1",
			Status:    ServerStatusRunning,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.NoError(t, s.Validate())
	})

	t.Run("missing id", func(t *testing.T) {
		s := Server{
			Name:      "Test",
			NodeID:    "node1",
			Status:    ServerStatusRunning,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, s.Validate())
	})

	t.Run("missing name", func(t *testing.T) {
		s := Server{
			ID:        "1",
			NodeID:    "node1",
			Status:    ServerStatusRunning,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, s.Validate())
	})

	t.Run("missing node_id", func(t *testing.T) {
		s := Server{
			ID:        "1",
			Name:      "Test",
			Status:    ServerStatusRunning,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, s.Validate())
	})

	t.Run("invalid status", func(t *testing.T) {
		s := Server{
			ID:        "1",
			Name:      "Test",
			NodeID:    "node1",
			Status:    ServerStatus("bogus"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, s.Validate())
	})

	t.Run("zero created_at", func(t *testing.T) {
		s := Server{
			ID:        "1",
			Name:      "Test",
			NodeID:    "node1",
			Status:    ServerStatusRunning,
			CreatedAt: time.Time{},
			UpdatedAt: now,
		}
		assert.Error(t, s.Validate())
	})

	t.Run("zero updated_at", func(t *testing.T) {
		s := Server{
			ID:        "1",
			Name:      "Test",
			NodeID:    "node1",
			Status:    ServerStatusRunning,
			CreatedAt: now,
			UpdatedAt: time.Time{},
		}
		assert.Error(t, s.Validate())
	})
}

func TestNodeModel(t *testing.T) {
	now := time.Now()

	node := Node{
		ID:        "node1",
		Name:      "Main Node",
		IPAddress: "192.168.1.1",
		Status:    NodeStatusOnline,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "node1", node.ID)
	assert.Equal(t, "Main Node", node.Name)
	assert.Equal(t, "192.168.1.1", node.IPAddress)
	assert.Equal(t, NodeStatusOnline, node.Status)
	assert.Equal(t, now, node.CreatedAt)
	assert.Equal(t, now, node.UpdatedAt)
}

func TestNodeStatus(t *testing.T) {
	testCases := []struct {
		status      NodeStatus
		stringValue string
	}{
		{NodeStatusOnline, "online"},
		{NodeStatusOffline, "offline"},
		{NodeStatusMaintenance, "maintenance"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.stringValue, string(tc.status))
	}
}

func TestValidNodeStatus(t *testing.T) {
	assert.True(t, ValidNodeStatus(NodeStatusOnline))
	assert.True(t, ValidNodeStatus(NodeStatusOffline))
	assert.True(t, ValidNodeStatus(NodeStatusMaintenance))
	assert.False(t, ValidNodeStatus(NodeStatus("invalid")))
	assert.False(t, ValidNodeStatus(NodeStatus("")))
}

func TestNodeValidate(t *testing.T) {
	now := time.Now()

	t.Run("valid node", func(t *testing.T) {
		n := Node{
			ID:        "node1",
			Name:      "Main",
			IPAddress: "10.0.0.1",
			Status:    NodeStatusOnline,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.NoError(t, n.Validate())
	})

	t.Run("missing id", func(t *testing.T) {
		n := Node{
			Name:      "Main",
			IPAddress: "10.0.0.1",
			Status:    NodeStatusOnline,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, n.Validate())
	})

	t.Run("missing ip_address", func(t *testing.T) {
		n := Node{
			ID:        "node1",
			Name:      "Main",
			Status:    NodeStatusOnline,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, n.Validate())
	})

	t.Run("invalid ip_address", func(t *testing.T) {
		n := Node{
			ID:        "node1",
			Name:      "Main",
			IPAddress: "not-an-ip",
			Status:    NodeStatusOnline,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, n.Validate())
	})

	t.Run("invalid status", func(t *testing.T) {
		n := Node{
			ID:        "node1",
			Name:      "Main",
			IPAddress: "10.0.0.1",
			Status:    NodeStatus("bogus"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, n.Validate())
	})
}

func TestBackupModel(t *testing.T) {
	now := time.Now()

	backup := Backup{
		ID:        "backup1",
		ServerID:  "server1",
		Path:      "/backups/server1/backup.zip",
		SizeBytes: 1024,
		Checksum:  "abc123",
		Status:    BackupStatusCompleted,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "backup1", backup.ID)
	assert.Equal(t, "server1", backup.ServerID)
	assert.Equal(t, "/backups/server1/backup.zip", backup.Path)
	assert.Equal(t, int64(1024), backup.SizeBytes)
	assert.Equal(t, "abc123", backup.Checksum)
	assert.Equal(t, BackupStatusCompleted, backup.Status)
	assert.Equal(t, now, backup.CreatedAt)
	assert.Equal(t, now, backup.UpdatedAt)
}

func TestBackupStatus(t *testing.T) {
	testCases := []struct {
		status      BackupStatus
		stringValue string
	}{
		{BackupStatusPending, "pending"},
		{BackupStatusRunning, "running"},
		{BackupStatusCompleted, "completed"},
		{BackupStatusFailed, "failed"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.stringValue, string(tc.status))
	}
}

func TestValidBackupStatus(t *testing.T) {
	assert.True(t, ValidBackupStatus(BackupStatusPending))
	assert.True(t, ValidBackupStatus(BackupStatusRunning))
	assert.True(t, ValidBackupStatus(BackupStatusCompleted))
	assert.True(t, ValidBackupStatus(BackupStatusFailed))
	assert.False(t, ValidBackupStatus(BackupStatus("invalid")))
	assert.False(t, ValidBackupStatus(BackupStatus("")))
}

func TestBackupValidate(t *testing.T) {
	now := time.Now()

	t.Run("valid backup", func(t *testing.T) {
		b := Backup{
			ID:        "b1",
			ServerID:  "s1",
			Path:      "/path",
			Status:    BackupStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.NoError(t, b.Validate())
	})

	t.Run("missing id", func(t *testing.T) {
		b := Backup{
			ServerID:  "s1",
			Path:      "/path",
			Status:    BackupStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, b.Validate())
	})

	t.Run("missing server_id", func(t *testing.T) {
		b := Backup{
			ID:        "b1",
			Path:      "/path",
			Status:    BackupStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, b.Validate())
	})

	t.Run("missing path", func(t *testing.T) {
		b := Backup{
			ID:        "b1",
			ServerID:  "s1",
			Status:    BackupStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, b.Validate())
	})

	t.Run("invalid status", func(t *testing.T) {
		b := Backup{
			ID:        "b1",
			ServerID:  "s1",
			Path:      "/path",
			Status:    BackupStatus("bogus"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Error(t, b.Validate())
	})
}
