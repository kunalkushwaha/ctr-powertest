package cmd

import (
	"fmt"
	"os"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/spf13/cobra"

	_ "github.com/containerd/containerd/differ"
	_ "github.com/containerd/containerd/linux"
	_ "github.com/containerd/containerd/metrics/cgroups"
	_ "github.com/containerd/containerd/services/containers"
	_ "github.com/containerd/containerd/services/content"
	_ "github.com/containerd/containerd/services/diff"
	_ "github.com/containerd/containerd/services/events"
	_ "github.com/containerd/containerd/services/healthcheck"
	_ "github.com/containerd/containerd/services/images"
	_ "github.com/containerd/containerd/services/namespaces"
	_ "github.com/containerd/containerd/services/snapshot"
	_ "github.com/containerd/containerd/services/tasks"
	_ "github.com/containerd/containerd/services/version"
	_ "github.com/containerd/containerd/snapshot/overlay"
)

var cfgFile string
var ctrRuntime libruntime.Runtime

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ctr-powertest",
	Short: "container runtime testing tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		//Run tests with new  server instance

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	RootCmd.PersistentFlags().StringP("runtime", "r", "containerd", "runtime [ containerd|crio ]")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
}
