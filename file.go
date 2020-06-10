package gnet

import (
	"github.com/slclub/link"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func SaveUploadFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// Read all  contents from file.
func ReadAll(file string) string {
	f, err := os.Open(file)
	if err != nil {
		link.ERROR("read file fail FILE:", file, "error:", err)
		return ""
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		link.ERROR("read file fail,ioutil.ReadAll FILE:", file, "error:", err)
		return ""
	}

	return string(fd)
}

// check the path is floder.
func IsDir(path string) bool {

	s, err := os.Stat(path)

	if err != nil {

		return false

	}

	return s.IsDir()
}
