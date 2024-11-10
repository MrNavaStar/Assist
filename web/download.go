package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/mrnavastar/assist/fs"
)

// Download downloads a file at the given path from a url if the file does not exist
func Download(file string, url string) (err error) {
	if fs.Exists(file) {
		return nil
	}
	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(file)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		os.Remove(file)
		return fmt.Errorf("failed to download: %s bad status: %s", url, resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(file)
		return err
	}
	return
}
