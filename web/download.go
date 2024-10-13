package web

import (
	"fmt"
	"github.com/mrnavastar/assist/fs"
	"io"
	"net/http"
	"os"
)

// Download downloads a file at the given path with the given name from a url if the file does not exist
func Download(filepath string, filename string, url string) (err error) {
	if fs.Exists(filepath + "/" + filename) {
		return nil
	}
	if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Create the file
	out, err := os.Create(filepath + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s bad status: %s", url, resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return
}
