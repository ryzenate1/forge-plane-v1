package downloadextract

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gamepanel/beacon/internal/installer/operations"
)

type DownloadExtract struct {
	URL     string `json:"url"`
	Dest    string `json:"dest"`
	Strip   int    `json:"strip,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

func init() {
	operations.Register("downloadExtract", factory)
}

func factory(args json.RawMessage) (operations.Operation, error) {
	var op DownloadExtract
	if err := json.Unmarshal(args, &op); err != nil {
		return nil, fmt.Errorf("downloadExtract: %w", err)
	}
	if op.URL == "" || op.Dest == "" {
		return nil, fmt.Errorf("downloadExtract: url and dest are required")
	}
	return &op, nil
}

func (op *DownloadExtract) Execute(ctx context.Context, serverDir string) error {
	dest := operations.ResolvePath(serverDir, op.Dest)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("create dest dir %q: %w", dest, err)
	}

	timeout := time.Duration(op.Timeout) * time.Second
	if op.Timeout <= 0 {
		timeout = operations.DefaultHTTPTimeout
	}

	resp, err := operations.HTTPGetWithRetry(ctx, op.URL, timeout)
	if err != nil {
		return fmt.Errorf("download %q: %w", op.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get %q: unexpected status %d", op.URL, resp.StatusCode)
	}

	log.Printf("[downloadExtract] downloading %s", op.URL)
	progress := &operations.DownloadProgress{Label: op.URL}
	teeReader := io.TeeReader(resp.Body, progress)

	contentType := resp.Header.Get("Content-Type")
	disposition := resp.Header.Get("Content-Disposition")

	if strings.Contains(contentType, "zip") || strings.HasSuffix(op.URL, ".zip") || strings.Contains(disposition, ".zip") {
		return op.extractZip(ctx, teeReader, dest)
	}
	if strings.Contains(contentType, "gzip") || strings.HasSuffix(op.URL, ".tar.gz") || strings.HasSuffix(op.URL, ".tgz") || strings.Contains(disposition, ".tar.gz") {
		return op.extractTarGz(ctx, teeReader, dest)
	}
	if strings.HasSuffix(op.URL, ".tar") || strings.Contains(contentType, "tar") {
		return op.extractTar(ctx, teeReader, dest)
	}
	return fmt.Errorf("unsupported archive format for %q (content-type: %s)", op.URL, contentType)
}

func (op *DownloadExtract) extractZip(ctx context.Context, r io.Reader, dest string) error {
	tmp, err := os.CreateTemp("", "gamepanel-dl-*.zip")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return fmt.Errorf("copy to temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}

	zipReader, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		name := op.stripPath(f.Name)
		if name == "" {
			continue
		}

		target := filepath.Join(dest, name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			continue
		}

		mode := f.Mode()
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, mode); err != nil {
				return fmt.Errorf("mkdir %q: %w", target, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), mode); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip entry %q: %w", f.Name, err)
		}

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
		if err != nil {
			rc.Close()
			return fmt.Errorf("create %q: %w", target, err)
		}

		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return fmt.Errorf("write %q: %w", target, err)
		}
	}
	return nil
}

func (op *DownloadExtract) extractTarGz(ctx context.Context, r io.Reader, dest string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gzr.Close()
	return op.extractTar(ctx, gzr, dest)
}

func (op *DownloadExtract) extractTar(ctx context.Context, r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read: %w", err)
		}

		name := op.stripPath(header.Name)
		if name == "" {
			continue
		}

		target := filepath.Join(dest, name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			continue
		}

		mode := header.FileInfo().Mode()

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, mode); err != nil {
				return fmt.Errorf("mkdir %q: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
			if err != nil {
				return fmt.Errorf("create %q: %w", target, err)
			}
			_, err = io.Copy(out, tr)
			out.Close()
			if err != nil {
				return fmt.Errorf("write %q: %w", target, err)
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}
			os.Remove(target)
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("symlink %q -> %q: %w", target, header.Linkname, err)
			}
		}
	}
	return nil
}

func (op *DownloadExtract) stripPath(name string) string {
	if op.Strip <= 0 {
		return name
	}
	parts := strings.SplitN(name, "/", op.Strip+1)
	if len(parts) <= op.Strip {
		return ""
	}
	return parts[op.Strip]
}
