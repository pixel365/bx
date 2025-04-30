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
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version information",
		Run: func(cmd *cobra.Command, _ []string) {
			printVersion(cmd)
		},
	}

	cmd.Flags().BoolP("verbose", "", false, "enable verbose mode")

	return cmd
}

func printVersion(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		fmt.Printf(
			"Version: %s\nCommit: %s\nDate: %s\nGo: %s %s/%s\n",
			buildVersion,
			buildCommit,
			buildDate,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
		)
		return
	}

	fmt.Printf("%s\n", buildVersion)
}
