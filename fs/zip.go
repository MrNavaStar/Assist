package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(filename string, files []string) error {
	arc, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer arc.Close()

	writer := zip.NewWriter(arc)
	defer writer.Close()

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return err
		}

		info, err := in.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file
		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		w, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(w, in); err != nil {
			return err
		}
		in.Close()
	}
	return nil
}

// Decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
// Credits to https://golangcode.com/unzip-files-in-go/
func Decompress(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		zipFile, err := f.Open()
		if err != nil {
			return nil, err
		}

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("%s: illegal file path", path)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		filenames = append(filenames, path)
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(outFile, zipFile)
		if err != nil {
			return nil, err
		}
		zipFile.Close()
	}
	return filenames, nil
}
