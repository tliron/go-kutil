package version

import (
	"runtime"
	"runtime/debug"

	"github.com/tliron/kutil/terminal"
)

func Print() {
	if GitVersion != "" {
		terminal.Printf("version=%s\n", GitVersion)
	}
	// TODO: use buildInfo.Settings's "vcs.revision" instead of version.GitRevision?
	if GitRevision != "" {
		terminal.Printf("revision=%s\n", GitRevision)
	}
	if Timestamp != "" {
		terminal.Printf("timestamp=%s\n", Timestamp)
	}
	terminal.Printf("arch=%s\n", runtime.GOARCH)
	terminal.Printf("os=%s\n", runtime.GOOS)
	terminal.Printf("compiler=%s\n", runtime.Compiler)
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		terminal.Printf("compiler-version=%s\n", buildInfo.GoVersion)
	}
}
