package cobra

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/version"
)

func NewVersionCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Show the version of %s", name),
		Long:  fmt.Sprintf(`Shows the version of %s.`, name),
		Run: func(cmd *cobra.Command, args []string) {
			// Why not use the version from runtime/debug.ReadBuildInfo? See:
			// https://github.com/golang/go/issues/29228
			if version.GitVersion != "" {
				fmt.Fprintf(terminal.Stdout, "version=%s\n", version.GitVersion)
			}
			if version.GitRevision != "" {
				fmt.Fprintf(terminal.Stdout, "revision=%s\n", version.GitRevision)
			}
			if version.Timestamp != "" {
				fmt.Fprintf(terminal.Stdout, "timestamp=%s\n", version.Timestamp)
			}
			fmt.Fprintf(terminal.Stdout, "arch=%s\n", runtime.GOARCH)
			fmt.Fprintf(terminal.Stdout, "os=%s\n", runtime.GOOS)
			fmt.Fprintf(terminal.Stdout, "compiler=%s\n", runtime.Compiler)
		},
	}
}
