package util

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CopyFile copy file from src to dest
func CopyFile(src, dest string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		log.Printf("Unable to open file '%s'", src)
		return err
	}
	defer in.Close()
	parentdir := filepath.Dir(dest)

	exist, err := Exists(parentdir)
	if !exist {
		os.MkdirAll(parentdir, os.ModePerm)
	}

	out, err := os.Create(dest)
	if err != nil {
		log.Printf("Unable to create dest file '%s'", dest)
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		log.Printf("Unable to copy file")
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	return nil
}

// ListAllFiles Get file pathes from specific folder and sub folders
func ListAllFiles(dir string) ([]string, error) {
	files := []string{}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Unable to read dir '%s'", dir)
		return nil, err
	}
	for _, fileInfo := range fileInfos {
		currentfile := filepath.Join(dir, fileInfo.Name())
		if fileInfo.IsDir() {
			subfiles, _ := ListAllFiles(currentfile)
			for _, f := range subfiles {
				files = append(files, f)
			}
		} else {
			files = append(files, currentfile)
		}
	}
	return files, err
}

// CopyDir copy folder from src to dest
func CopyDir(src, dest string) (err error) {
	files, err := ListAllFiles(src)
	if err != nil {
		log.Printf("Unable list all files under '%s'", src)
		return err
	}
	for _, file := range files {
		relativepath, _ := filepath.Rel(src, file)
		_, dirname := filepath.Split(src)
		joinedDest := filepath.Join(dest, dirname, relativepath)
		CopyFile(file, joinedDest)
	}
	return nil
}

// Exists check file existance
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Unzip unzip zip file to dest folder
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		log.Printf("Unable to open reader to '%s'", src)
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			log.Printf("Unable to open zip file '%s'", src)
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Printf("Unable to open file '%s'", path)
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				log.Printf("Unable to copy file")
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			log.Printf("Unable to extract and file")
			return err
		}
	}

	return nil
}

// Replace replace string content in file
func Replace(filepath, oldstr, newstr string) error {
	read, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Unable to read file '%s'", filepath)
		return err
	}
	newContents := strings.Replace(string(read), oldstr, newstr, -1)
	err = ioutil.WriteFile(filepath, []byte(newContents), 0)
	if err != nil {
		log.Printf("Unable to write file '%s'", filepath)
		return err
	}
	return nil
}
