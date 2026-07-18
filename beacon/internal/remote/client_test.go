package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientRejectsNon2xxResponses(t *testing.T) {
	tests := []struct {
		name   string
		status int
		call   func(Client) error
	}{
		{
			name:   "get 401",
			status: http.StatusUnauthorized,
			call: func(client Client) error {
				_, err := client.GetServerConfiguration(context.Background(), "server-id")
				return err
			},
		},
		{
			name:   "post 404",
			status: http.StatusNotFound,
			call: func(client Client) error {
				return client.SendServerStats(context.Background(), "server-id", ServerStats{})
			},
		},
		{
			name:   "post 500",
			status: http.StatusInternalServerError,
			call: func(client Client) error {
				return client.ResetServersState(context.Background())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte("panel failure"))
			}))
			defer server.Close()

			err := test.call(NewClient(server.URL+"/api/v1", "token"))
			if err == nil {
				t.Fatal("expected non-2xx response to fail")
			}
			if !strings.Contains(err.Error(), fmt.Sprintf("%d %s", test.status, http.StatusText(test.status))) {
				t.Fatalf("expected status in error, got %v", err)
			}
			if !strings.Contains(err.Error(), "panel failure") {
				t.Fatalf("expected bounded response body in error, got %v", err)
			}
		})
	}
}

func TestGetChecksStatusBeforeDecoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("this is not json"))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token").GetServers(context.Background(), 10)
	if err == nil || !strings.Contains(err.Error(), "401 Unauthorized") {
		t.Fatalf("expected HTTP status error, got %v", err)
	}
	if strings.Contains(err.Error(), "decode") {
		t.Fatalf("expected status to be checked before decode, got %v", err)
	}
}

func TestClientRejectsMalformedSuccessJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token").GetServers(context.Background(), 10)
	if err == nil || !strings.Contains(err.Error(), "decode servers response") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestClientSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remote/servers" || r.URL.Query().Get("per_page") != "25" {
			t.Fatalf("unexpected request URL %s", r.URL.String())
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Fatalf("unexpected authorization header %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{"uuid": "123e4567-e89b-12d3-a456-426614174000"}},
		})
	}))
	defer server.Close()

	servers, err := NewClient(server.URL+"/api/remote", "token").GetServers(context.Background(), 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 || servers[0].Uuid != "123e4567-e89b-12d3-a456-426614174000" {
		t.Fatalf("unexpected servers response: %+v", servers)
	}
}

func TestActivityLogsUseForgeActionField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remote/activity" {
			t.Fatalf("unexpected request path %q", r.URL.Path)
		}
		var payload struct {
			Data []struct {
				Action string `json:"action"`
				Event  string `json:"event"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode activity payload: %v", err)
		}
		if len(payload.Data) != 1 || payload.Data[0].Action != "sftp.file.read" || payload.Data[0].Event != "" {
			t.Fatalf("unexpected activity payload: %+v", payload)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := NewClient(server.URL, "token").SendActivityLogs(context.Background(), []Activity{{Event: "sftp.file.read"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHeartbeatUsesAPIV1WhileRemoteCallsUseRemotePrefix(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL+"/api/remote", "token")
	if err := client.SendNodeHeartbeat(context.Background(), "node-id", NodeHeartbeat{}); err != nil {
		t.Fatal(err)
	}
	if err := client.ResetServersState(context.Background()); err != nil {
		t.Fatal(err)
	}

	want := []string{"/api/v1/nodes/node-id/heartbeat", "/api/remote/servers/reset"}
	if fmt.Sprint(paths) != fmt.Sprint(want) {
		t.Fatalf("unexpected routes: got %v want %v", paths, want)
	}
}

func TestNewClientURLNormalization(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantRemote string
		wantAPI    string
	}{
		{"bare URL", "http://panel.test", "http://panel.test/api/remote", "http://panel.test/api/v1"},
		{"trailing slash", "http://panel.test/", "http://panel.test/api/remote", "http://panel.test/api/v1"},
		{"whitespace", "  http://panel.test  ", "http://panel.test/api/remote", "http://panel.test/api/v1"},
		{"already /api/remote", "http://panel.test/api/remote", "http://panel.test/api/remote", "http://panel.test/api/v1"},
		{"already /api/v1", "http://panel.test/api/v1", "http://panel.test/api/remote", "http://panel.test/api/v1"},
		{"with port", "http://panel.test:8080", "http://panel.test:8080/api/remote", "http://panel.test:8080/api/v1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewClient(tc.input, "test-token")
			cl := c.(*client)
			if cl.remoteBaseURL != tc.wantRemote {
				t.Errorf("remoteBaseURL = %q, want %q", cl.remoteBaseURL, tc.wantRemote)
			}
			if cl.apiBaseURL != tc.wantAPI {
				t.Errorf("apiBaseURL = %q, want %q", cl.apiBaseURL, tc.wantAPI)
			}
		})
	}
}

func TestGetServerConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remote/servers/test-uuid" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"settings":{},"process_configuration":{"startup":{"done":["ok"],"strip_ansi":false},"stop":{"type":"command","value":"stop"}},"mounts":[]}`))
	}))
	defer server.Close()

	cfg, err := NewClient(server.URL, "token").GetServerConfiguration(context.Background(), "test-uuid")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ProcessConfiguration == nil {
		t.Fatal("expected process configuration")
	}
	if cfg.ProcessConfiguration.Stop.Type != "command" {
		t.Fatalf("stop type = %q", cfg.ProcessConfiguration.Stop.Type)
	}
}

func TestCreatePlacementReservation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remote/reservations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"res-1","nodeId":"node-1","reservationType":"create","cpu":2,"memory":4096,"disk":10240,"status":"pending","expiresAt":"2026-01-01T00:00:00Z","createdAt":"2026-01-01T00:00:00Z"}`))
	}))
	defer server.Close()

	res, err := NewClient(server.URL, "token").CreatePlacementReservation(context.Background(), PlacementReservationRequest{
		NodeID:          "node-1",
		ReservationType: "create",
		CPU:             2,
		Memory:          4096,
		Disk:            10240,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ID != "res-1" || res.Status != "pending" {
		t.Fatalf("unexpected reservation: %+v", res)
	}
}

func TestCreatePlacementReservationDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token").CreatePlacementReservation(context.Background(), PlacementReservationRequest{})
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func TestClientContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewClient(server.URL, "token").GetServers(ctx, 10)
	if err == nil {
		t.Fatal("expected error with cancelled context")
	}
}

func TestClientAuthHeaders(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	_ = NewClient(server.URL, "my-secret-token").ResetServersState(context.Background())
	if authHeader != "Bearer my-secret-token" {
		t.Fatalf("Authorization header = %q, want %q", authHeader, "Bearer my-secret-token")
	}
}
