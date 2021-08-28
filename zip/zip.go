package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(source, target string, verbose bool) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = strings.Replace(path, "./", "", -1)
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		if err == nil && verbose {
			fmt.Println(path)
		}
		return err
	})
	return err
}

func Unzip(src string, dest string, verbose, recursive, deleteSrc bool) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// ignore .DS_Store
		if strings.HasSuffix(strings.ToLower(f.Name), ".ds_store") {
			continue
		}
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		if verbose {
			fmt.Println(fpath)
		}
		filenames = append(filenames, fpath)
		if err := doUnzip(f, fpath); err != nil {
			return filenames, err
		}
		if recursive && strings.HasSuffix(fpath, "zip") {
			if _, err := Unzip(fpath, filepath.Dir(fpath), verbose, recursive, deleteSrc); err != nil {
				return filenames, err
			}
		}
	}
	if deleteSrc && strings.HasSuffix(src, "zip") {
		os.Remove(src)
	}
	return filenames, nil
}

func doUnzip(f *zip.File, fpath string) error {
	if f.FileInfo().IsDir() {
		// Make Folder
		return os.MkdirAll(fpath, os.ModePerm)
	}

	// Make File
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)

	if err != nil {
		return err
	}
	return nil
}
