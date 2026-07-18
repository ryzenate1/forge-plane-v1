package operations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ResolvePath(serverDir, target string) string {
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}
	return filepath.Clean(filepath.Join(serverDir, target))
}

func FileExists(serverDir, target string) (bool, error) {
	full := ResolvePath(serverDir, target)
	info, err := os.Stat(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("stat %q: %w", full, err)
	}
	return info.Mode().IsRegular(), nil
}

func EnsureParentDir(path string) error {
	parent := filepath.Dir(path)
	if parent == "." {
		return nil
	}
	return os.MkdirAll(parent, 0o755)
}

func PathExists(serverDir, target string) (bool, error) {
	full := ResolvePath(serverDir, target)
	_, err := os.Stat(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("stat %q: %w", full, err)
	}
	return true, nil
}

const (
	DefaultRetries    = 3
	DefaultHTTPTimeout = 10 * time.Minute
	MetaHTTPTimeout   = 30 * time.Second
)

func retryBackoff(attempt int) time.Duration {
	return time.Duration(attempt) * time.Second
}

type DownloadProgress struct {
	Label string
	total int64
	last  time.Time
}

func (dp *DownloadProgress) Write(p []byte) (int, error) {
	n := len(p)
	dp.total += int64(n)
	if time.Since(dp.last) > 5*time.Second {
		log.Printf("[%s] %d bytes downloaded", dp.Label, dp.total)
		dp.last = time.Now()
	}
	return n, nil
}

func HTTPGetWithRetry(ctx context.Context, url string, timeout time.Duration) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= DefaultRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryBackoff(attempt)):
			}
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("User-Agent", "GamePanel-Beacon/1.0")
		client := &http.Client{Timeout: timeout}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("request failed after %d retries: %w", DefaultRetries, lastErr)
}

func DownloadWithRetry(ctx context.Context, url, dest string, timeout time.Duration) (int64, error) {
	resp, err := HTTPGetWithRetry(ctx, url, timeout)
	if err != nil {
		return 0, fmt.Errorf("download %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GET %q: status %d", url, resp.StatusCode)
	}

	if err := EnsureParentDir(dest); err != nil {
		return 0, fmt.Errorf("create parent dir: %w", err)
	}

	out, err := os.Create(dest)
	if err != nil {
		return 0, fmt.Errorf("create %q: %w", dest, err)
	}
	defer out.Close()

	progress := &DownloadProgress{Label: dest}
	writer := io.MultiWriter(out, progress)
	written, err := io.Copy(writer, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("write %q: %w", dest, err)
	}
	if written == 0 {
		return 0, fmt.Errorf("downloaded file %q is empty", dest)
	}
	return written, nil
}

func VerifySHA256(path, expectedHex string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open for checksum: %w", err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("checksum read: %w", err)
	}
	got := hex.EncodeToString(h.Sum(nil))
	if got != expectedHex {
		return fmt.Errorf("checksum mismatch: got %s, expected %s", got, expectedHex)
	}
	return nil
}
