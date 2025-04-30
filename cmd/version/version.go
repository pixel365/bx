package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version information",
		Run: func(cmd *cobra.Command, _ []string) {
			printVersion()
		},
	}
}

func printVersion() {
	fmt.Printf(
		"Version: %s\nCommit: %s\nDate: %s\nGo: %s %s/%s\n",
		buildVersion,
		buildCommit,
		buildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
