package cmd

import (
	"context"

	"github.com/kunalkushwaha/ctr-powertest/libruntime"
	"github.com/kunalkushwaha/ctr-powertest/testcase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// basicCmd represents the basic command
var basicCmd = &cobra.Command{
	Use:   "basic",
	Short: "runs basic tests",
	Run:   runBasicTest,
}

var config = libruntime.RuntimeConfig{
	RuntimeName:      "containerd",
	RunDefaultServer: true,
	Root:             "/var/lib/powertest", //containerd-ctr-powertest
	RuntimeEndpoint:  "/run/powertest/containerd.sock",
	DebugEndpoint:    "/run/powertest/debug.sock",
}

func init() {
	RootCmd.AddCommand(basicCmd)
	log.SetLevel(log.InfoLevel)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// basicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// basicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func runBasicTest(cmd *cobra.Command, args []string) {

	testSetup, err := testcase.SetupTestEnvironment(config.RuntimeName, config, false)
	if err != nil {
		log.Fatal("Error while setting up environment : ", err)
	}

	err = testSetup.RunAllTests(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
