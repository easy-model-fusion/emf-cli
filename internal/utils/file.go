package utils

import (
	"embed"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
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
