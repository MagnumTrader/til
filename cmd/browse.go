package cmd

import (

	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "open up neovim in the til folder",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEditor()
	},
}
