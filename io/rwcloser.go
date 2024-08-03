package io

import "io"

// Merges a ReadCloser and a WriteCloser into one
type RWCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (rw RWCloser) Close() error {
	err := rw.ReadCloser.Close()
	if err := rw.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}
