package main

import (
	_ "github.com/tliron/kutil/ard"
	_ "github.com/tliron/kutil/cobra"
	_ "github.com/tliron/kutil/exec"
	_ "github.com/tliron/kutil/fswatch"
	_ "github.com/tliron/kutil/js"
	_ "github.com/tliron/kutil/kubernetes"
	_ "github.com/tliron/kutil/logging"
	_ "github.com/tliron/kutil/logging/journal"
	_ "github.com/tliron/kutil/logging/klog"
	_ "github.com/tliron/kutil/logging/simple"
	_ "github.com/tliron/kutil/logging/sink"
	_ "github.com/tliron/kutil/logging/zerolog"
	_ "github.com/tliron/kutil/problems"
	_ "github.com/tliron/kutil/protobuf"
	_ "github.com/tliron/kutil/reflection"
	_ "github.com/tliron/kutil/terminal"
	_ "github.com/tliron/kutil/transcribe"
	_ "github.com/tliron/kutil/url"
	_ "github.com/tliron/kutil/util"
	_ "github.com/tliron/kutil/version"
)

func main() {
}
