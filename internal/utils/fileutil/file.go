package fileutil

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"embed"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

const invalidFileNameCharsRegex = "\\\\|/|:|\\*|\\?|<|>"

// IsFileNameValid returns true if the name is a valid file name.
func IsFileNameValid(name string) bool {
	re := regexp.MustCompile(invalidFileNameCharsRegex)
	return !re.MatchString(name)
}

// ValidFileName returns an error if the which arg is not a valid file name.
func ValidFileName(which int, optional bool) cobra.PositionalArgs {
	if which <= 0 {
		panic("which must be strictly positive")
	}

	which -= 1 // real index is 0-based

	return func(cmd *cobra.Command, args []string) error {
		if len(args) <= which {
			if optional {
				return nil
			}
			return fmt.Errorf("requires at least %d arg(s), only received %d", which+1, len(args))
		}

		name := args[which]
		if !IsFileNameValid(name) {
			return fmt.Errorf("'%s' is not a valid file name", name)
		}

		return nil
	}
}

// CopyEmbeddedFile copies an embedded file to a destination.
func CopyEmbeddedFile(fs embed.FS, file, dst string) error {
	content, err := fs.ReadFile(file)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, content, os.ModePerm)
}

// CloseFile closes a file and logs an error if it occurs.
func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		pterm.Error.Println(fmt.Sprintf("Error closing file: %s", err))
	}
}

// IsExistingPath check if the requested path exists
func IsExistingPath(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil && !os.IsNotExist(err) {
		// An error occurred while verifying the non-existence of the path
		pterm.Error.Println(fmt.Sprintf("Error checking the existence of %s : %s", name, err))
		return false, err
	} else if err == nil {
		// Path already exists
		return true, nil
	}
	// Path does not exist
	return false, nil
}

// IsDirectoryEmpty check if the requested directory is empty
func IsDirectoryEmpty(name string) (bool, error) {
	file, err := os.Open(name)
	defer CloseFile(file)
	if err != nil {
		return false, err
	}

	_, err = file.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// DeleteDirectoryIfEmpty checks if the specified directory is empty. If it is empty, the directory is removed.
func DeleteDirectoryIfEmpty(path string) error {
	// Check if the directory is empty
	if empty, err := IsDirectoryEmpty(path); err != nil {
		return err
	} else if empty {
		// Current directory is empty : removing it
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// MoveFiles moves all files from the source directory to the destination directory
func MoveFiles(sourceDir, destinationDir string) error {
	fileList, err := filepath.Glob(filepath.Join(sourceDir, "*"))
	if err != nil {
		return err
	}

	for _, file := range fileList {
		_, fileName := filepath.Split(file)
		destinationPath := filepath.Join(destinationDir, fileName)

		err = os.Rename(file, destinationPath)
		if err != nil {
			return err
		}

	}

	return nil
}

// Compress compresses a file or a folder
// Source: https://github.com/mimoo/eureka/blob/master/folders.go
func Compress(src, dst string) error {
	// is file a folder and does it exist?
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// create a new file dst
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer CloseFile(out)

	// create a new buffer writer
	buf := bufio.NewWriter(out)
	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	mode := fi.Mode()
	if mode.IsRegular() {
		// get header
		header, err := tar.FileInfoHeader(fi, src)
		if err != nil {
			return err
		}
		// write header
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// get content
		data, err := os.Open(src)
		if err != nil {
			return err
		}
		if _, err = io.Copy(tw, data); err != nil {
			return err
		}
	} else if mode.IsDir() { // folder

		// walk through every file in the folder
		err = filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			// generate tar header
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// must provide real name
			// (see https://golang.org/src/archive/tar/common.go?#L626)
			header.Name = filepath.ToSlash(file)

			// write header
			if err = tw.WriteHeader(header); err != nil {
				return err
			}
			// if not a dir, write file content
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err = io.Copy(tw, data); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("error: file type not supported")
	}

	// produce tar
	if err = tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err = zr.Close(); err != nil {
		return err
	}

	err = buf.Flush()
	if err != nil {
		return err
	}
	return nil
}
