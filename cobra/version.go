package cobra

import (
	"fmt"
	"runtime"
	"runtime/debug"

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
				terminal.Printf("version=%s\n", version.GitVersion)
			}
			if version.GitRevision != "" {
				terminal.Printf("revision=%s\n", version.GitRevision)
			}
			if version.Timestamp != "" {
				terminal.Printf("timestamp=%s\n", version.Timestamp)
			}
			terminal.Printf("arch=%s\n", runtime.GOARCH)
			terminal.Printf("os=%s\n", runtime.GOOS)
			terminal.Printf("compiler=%s\n", runtime.Compiler)
			if buildInfo, ok := debug.ReadBuildInfo(); ok {
				terminal.Printf("compiler-version=%s\n", buildInfo.GoVersion)
			}
		},
	}
}
