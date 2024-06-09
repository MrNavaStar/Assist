package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func compress(filename string, files []string) error {
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
		defer in.Close()

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
	}
	return nil
}

// decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
// Credits to https://golangcode.com/unzip-files-in-go/
func decompress(src string, dest string) ([]string, error) {
	var outFile *os.File
	var zipFile io.ReadCloser
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	clean := func() {
		if outFile != nil {
			outFile.Close()
			outFile = nil
		}

		if zipFile != nil {
			zipFile.Close()
			zipFile = nil
		}
	}

	for _, f := range r.File {
		zipFile, err = f.Open()
		if err != nil {
			return filenames, err
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			clean()
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			clean()
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			clean()
			return filenames, err
		}

		outFile, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			clean()
			return filenames, err
		}

		_, err = io.Copy(outFile, zipFile)
		clean()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
