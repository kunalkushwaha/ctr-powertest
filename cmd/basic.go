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
	Root:             "/var/lib/powertest",
	RuntimeEndpoint:  "/run/powertest/containerd.sock",
	DebugEndpoint:    "/run/powertest/debug.sock",
}

var stdConfig = libruntime.RuntimeConfig{
	RuntimeName:      "containerd",
	RunDefaultServer: false,
	Root:             "/var/lib/containerd",
	RuntimeEndpoint:  "/run/containerd/containerd.sock",
	DebugEndpoint:    "/run/containerd/debug.sock",
}

func init() {
	RootCmd.AddCommand(basicCmd)

	//Set Log level
	log.SetLevel(log.DebugLevel)

}

func runBasicTest(cmd *cobra.Command, args []string) {

	//Run tests with new  server instance
	ctrRuntime, err := testcase.SetupTestEnvironment(stdConfig.RuntimeName, stdConfig, false)
	if err != nil {
		log.Fatal("Error while setting up environment : ", err)
	}

	var singleClientTestCases testcase.Testcases

	singleClientTestCases = &testcase.BasicContainerTest{Runtime: ctrRuntime}

	err = singleClientTestCases.RunAllTests(context.TODO(), args)
	if err != nil {
		log.Fatal(err)
	}
}
