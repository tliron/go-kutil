package cobra

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func SetFlagsFromEnvironment(prefix string, command *cobra.Command) {
	setFlagsFromEnvironment(prefix, command.PersistentFlags())
	setFlagsFromEnvironment(prefix, command.Flags())
}

func setFlagsFromEnvironment(prefix string, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if value, ok := os.LookupEnv(prefix + flag.Name); ok {
			flags.Set(flag.Name, value)
		}
	})
}
