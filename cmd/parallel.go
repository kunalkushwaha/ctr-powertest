package cmd

import (
	"github.com/kunalkushwaha/ctr-powertest/testcase"
	//	"context"
	//	"log"

	"context"
	"log"

	"github.com/spf13/cobra"
)

// parallelCmd represents the parallel command
var parallelCmd = &cobra.Command{
	Use:   "parallel",
	Short: "Run container tests in parallel (Stress Test)",
	Run:   runParallelTest,
}

func init() {
	RootCmd.AddCommand(parallelCmd)
}

func runParallelTest(cmd *cobra.Command, args []string) {

	//Run tests with new  server instance
	ctrRuntime, err := testcase.SetupTestEnvironment(stdConfig.RuntimeName, stdConfig, false)
	if err != nil {
		log.Fatal("Error while setting up environment : ", err)
	}

	var parallelTestCases testcase.Testcases
	parallelTestCases = &testcase.ParallelContainerTest{Runtime: ctrRuntime}
	err = parallelTestCases.RunAllTests(context.TODO(), args)
	if err != nil {
		log.Fatal(err)
	}

}
