package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	expectedArgs = 2
	dirPerm = 0755
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

// TODO: extract args to a config file that users can edit
// NOTE: Used to ignore my custom session handler for neovim

var editorargs = []string{"--cmd", "lua vim.g.nosession=1"}

var rootCmd = &cobra.Command{
	Use:   "til <Path> <Title>",
	Short: "Today I Learned CLI tool",
	Long:  "til is a tool to register til posts in a specified folder",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != expectedArgs {
			return fmt.Errorf("til expects %d arguments <Path> and <Title>", expectedArgs)
		}

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

		editorargs = append(editorargs, filePath)

		nvim := exec.Command(editor, editorargs...)

		nvim.Stdin = os.Stdin
		nvim.Stdout = os.Stdout
		nvim.Stderr = os.Stderr

		return nvim.Run()
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
