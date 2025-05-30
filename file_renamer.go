package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type fileRenamer struct {
	fileSystem billy.Filesystem
	stderr     io.Writer
}

func newFileRenamer(fs billy.Filesystem, stderr io.Writer) *fileRenamer {
	return &fileRenamer{fs, stderr}
}

func (r *fileRenamer) Rename(tr textRenamer, path string, baseDir string, verbose bool) error {
	p, err := r.renamePath(tr, path, baseDir)
	if err != nil {
		return err
	}

	if p != path {
		if verbose {
			err := r.print("Moving", path)
			if err != nil {
				return err
			}
		}

		err = r.fileSystem.Rename(path, p)
		if err != nil {
			return err
		}
	}

	ok, err := r.isTextFile(p)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	f, err := r.fileSystem.OpenFile(p, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	bs, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	bbs := []byte(tr.Rename(string(bs)))
	if bytes.Equal(bs, bbs) {
		return nil
	}

	if verbose {
		err := r.print("Writing", path)
		if err != nil {
			return err
		}
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(bbs)
	return err
}

func (r *fileRenamer) isTextFile(path string) (bool, error) {
	f, err := r.fileSystem.Open(path)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	bs := make([]byte, fileTypeDetectionBufferSize)
	_, err = f.Read(bs)
	if err != nil && err != io.EOF {
		return false, err
	}

	t, err := filetype.Match(bs)
	if err != nil {
		return false, err
	}

	return t == types.Unknown, nil
}

func (r *fileRenamer) renamePath(tr textRenamer, path, baseDir string) (string, error) {
	if baseDir == "" {
		return path, nil
	}

	b, err := filepath.Rel(baseDir, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, tr.Rename(b)), nil
}

func (r *fileRenamer) print(xs ...interface{}) error {
	_, err := fmt.Fprintln(r.stderr, xs...)
	return err
}
