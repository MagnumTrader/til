package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// TODO: Add edit functionality from til
// TODO: Add autocompletion
// TODO: Add a ls command that prints a tree

const (
	dirPerm      = 0755
	editor       = "nvim"
	fileHeader   = "# What did you learn today"
	fileExt      = ".md"
)

var tilDir = filepath.Join(must(os.UserHomeDir()), "til")

func must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func init() {
	rootCmd.AddCommand(browseCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "til <Path> <Title>",
	Short: "Today I Learned CLI tool",
	Long:  "til is a tool to register til posts in a specified folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tilPath := args[0]
		title := strings.ReplaceAll(args[1], " ", "_")

		folderpath := filepath.Join(tilDir, tilPath)
		os.MkdirAll(folderpath, dirPerm)

		filePath := fmt.Sprintf("%s/%s%s", folderpath, title, fileExt)

		file, err := os.Create(filePath)
		if err != nil {
			return err
		}

		fields := map[string]string{
			"date": time.Now().Format(time.DateTime),
		}
		err = WriteFrontMatter(file, fields)

		if err != nil {
			return err
		}

		file.WriteString(fileHeader)

		return runEditor(filePath)
	},
}

func WriteFrontMatter(f *os.File, fields map[string]string) error {

	_, err := f.WriteString("---\n")
	if err != nil {
		return err
	}

	for key, value := range fields {
		_, err := fmt.Fprintf(f, "%s: %s\n", key, value)
		if err != nil {
			return err
		}

	}
	_, err = f.WriteString("---\n")
	return err
}

// Will run the editor in the til directory with the provided arguments.
func runEditor(args ...string) error {
	nvim := exec.Command(editor, args...)

	nvim.Stdin = os.Stdin
	nvim.Stdout = os.Stdout
	nvim.Stderr = os.Stderr
	nvim.Dir = tilDir

	return nvim.Run()
}
