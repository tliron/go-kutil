package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
)

// TODO: This is not needed anymore with new versions of Cobra

func NewBashCompletionCommand(name string, rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: fmt.Sprintf("Generate bash completion script for %s", name),
		Long:  fmt.Sprintf(`Generates bash completion script for %s.`, name),
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(terminal.Stdout)
		},
	}
}
