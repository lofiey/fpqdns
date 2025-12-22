package geo

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Downloader struct {
	GeoIPFile   string
	GeoSiteFile string

	GeoIPURL   string
	GeoSiteURL string
}

func (d *Downloader) Ensure() error {
	if err := ensureFile(d.GeoIPFile, d.GeoIPURL); err != nil {
		return err
	}
	if err := ensureFile(d.GeoSiteFile, d.GeoSiteURL); err != nil {
		return err
	}
	return nil
}

func ensureFile(path, url string) error {
	if path == "" || url == "" {
		return nil
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("geo download failed: " + resp.Status)
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
