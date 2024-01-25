package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func IsTarGz(name string) bool {
	buffer, err := readData(name)
	if err != nil && err != io.EOF {
		return false
	}

	if err != nil && err != io.EOF {
		// return fmt.Errorf("file '%s' cannot be read: %s", name, err)
		return false
	}

	if contentType := http.DetectContentType(buffer); contentType != "application/x-gzip" {
		return false
	}
	return true
}

func IsZip(name string) bool {
	buffer, err := readData(name)
	if err != nil && err != io.EOF {
		return false
	}

	if contentType := http.DetectContentType(buffer); contentType != "application/zip" {
		return false
	}
	return true
}

func IsYaml(name string) bool {
	buffer, err := readData(name)
	if err != nil && err != io.EOF {
		return false
	}

	fmt.Println("test content type: ", http.DetectContentType(buffer))
	if contentType := http.DetectContentType(buffer); !strings.Contains(contentType, "text/plain") {
		return false
	}
	return true
}

func readData(name string) (buffer []byte, err error) {
	raw, err := os.Open(name)
	if err != nil {
		return
	}

	defer func() {
		_ = raw.Close()
		_, _ = raw.Seek(0, 0)
	}()
	buffer = make([]byte, 512)
	_, err = raw.Read(buffer)
	
	if err != nil {
		return
	}
	return
}

func ExtractZip(path, artifactPath string) error {
	zipReader, err := zip.OpenReader(artifactPath)
	defer func() {
		_ = zipReader.Close()
	}()

	if err != nil {
		fmt.Println("inside extract zip", err)
		return ErrExtractZip(err, path)
	}
	buffer := make([]byte, 1<<4)
	for _, file := range zipReader.File {
		fmt.Println("inside loop zip", file.Name)

		fd, err := file.Open()
		defer func() {
			_ = fd.Close()
		}()

		if err != nil {
			return ErrExtractZip(err, path)
		}

		filePath := filepath.Join(path, file.Name)

		if file.FileInfo().IsDir() {
			err := os.Mkdir(file.Name, file.Mode())
			if err != nil {
				return ErrExtractZip(err, path)
			}
		} else {
			openedFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return ErrExtractZip(err, path)
			}
			io.CopyBuffer(openedFile, fd, buffer)
			defer func() {
				_ = openedFile.Close()
			}()
		}

	}
	return nil

}

func ExtractTarGz(path string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return ErrExtractTarXZ(err, path)
	}
	defer func(){
		_ = uncompressedStream.Close()
	}()

	tarReader := tar.NewReader(uncompressedStream)
	fmt.Println("insinde extact ")

	if err != nil {
		return ErrExtractTarXZ(err, path)
	}

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return ErrExtractTarXZ(err, path)
		}
		fmt.Println("insinde extact loop: ", header.Typeflag-'0', header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(path, header.Name), 0755); err != nil {
				return ErrExtractTarXZ(err, path)
			}
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Join(path, filepath.Dir(header.Name)), 0755)
			outFile, err := os.Create(filepath.Join(path, header.Name))
			if err != nil {
				return ErrExtractTarXZ(err, path)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return ErrExtractTarXZ(err, path)
			}
			outFile.Close()

		default:
			return ErrExtractTarXZ(err, path)
		}

	}
	return nil
}

func ProcessExtractedContent(dirPath string, f func(path string) error) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ErrReadDir(err, dirPath)
	}

	for _, entry := range entries {
		err := f(filepath.Join(dirPath, entry.Name()))
		fmt.Println("INSIDE PROCESS >>>>>>>>>", entry.Name())
		if err != nil {
			return err
		}
	}
	return nil
}