package downloadfile

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gamepanel/beacon/internal/installer/operations"
)

type DownloadFile struct {
	URL      string `json:"url"`
	Dest     string `json:"dest"`
	Timeout  int    `json:"timeout,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Mode     int    `json:"mode,omitempty"`
}

func init() {
	operations.Register("downloadFile", factory)
}

func factory(args json.RawMessage) (operations.Operation, error) {
	var op DownloadFile
	if err := json.Unmarshal(args, &op); err != nil {
		return nil, fmt.Errorf("downloadFile: %w", err)
	}
	if op.URL == "" || op.Dest == "" {
		return nil, fmt.Errorf("downloadFile: url and dest are required")
	}
	return &op, nil
}

func (op *DownloadFile) Execute(ctx context.Context, serverDir string) error {
	dest := operations.ResolvePath(serverDir, op.Dest)

	timeout := time.Duration(op.Timeout) * time.Second
	if op.Timeout <= 0 {
		timeout = operations.DefaultHTTPTimeout
	}

	if _, err := operations.DownloadWithRetry(ctx, op.URL, dest, timeout); err != nil {
		return err
	}

	if op.Checksum != "" {
		if err := operations.VerifySHA256(dest, op.Checksum); err != nil {
			return fmt.Errorf("checksum: %w", err)
		}
	}

	if op.Mode > 0 {
		if err := os.Chmod(dest, os.FileMode(op.Mode)); err != nil {
			return fmt.Errorf("chmod %q: %w", dest, err)
		}
	}

	return nil
}
