package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/go-git/go-billy/v5"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type fileRenamer struct {
	fileSystem billy.Filesystem
}

func newFileRenamer(fs billy.Filesystem) *fileRenamer {
	return &fileRenamer{fs}
}

func (r *fileRenamer) Rename(tr *caseTextRenamer, path string) error {
	ok, err := r.validatePath(path)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	p := tr.Rename(path)

	if p != path {
		err := r.fileSystem.Rename(path, p)
		if err != nil {
			return err
		}
	}

	i, err := r.fileSystem.Lstat(p)
	if err != nil {
		return err
	} else if i.IsDir() {
		return nil
	}

	ok, err = r.isTextFile(p)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	f, err := r.fileSystem.OpenFile(p, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	bbs := []byte(tr.Rename(string(bs)))
	if bytes.Equal(bs, bbs) {
		return nil
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

func (*fileRenamer) validatePath(s string) (bool, error) {
	ok, err := regexp.MatchString("(^|/)\\.", s)
	if err != nil {
		return false, err
	}

	return !ok, nil
}

func (r *fileRenamer) isTextFile(path string) (bool, error) {
	f, err := r.fileSystem.Open(path)
	if err != nil {
		return false, err
	}

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
