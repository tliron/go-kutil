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
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if version.GitVersion != "" {
				terminal.Printf("version=%s\n", version.GitVersion)
			}
			// TODO: use buildInfo.Settings's "vcs.revision" instead of version.GitRevision?
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
				/*
					for _, setting := range buildInfo.Settings {
						terminal.Printf("compiler-setting=%s=%s\n", setting.Key, setting.Value)
					}
				*/
			}
		},
	}
}
