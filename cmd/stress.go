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
	Use:   "stress",
	Short: "Run container tests in parallel (Stress Test)",
	Run:   runStressTest,
}

func init() {
	RootCmd.AddCommand(parallelCmd)
}

func runStressTest(cmd *cobra.Command, args []string) {
	initTestSuite(cmd)

	var stressTestCases testcase.Testcases
	stressTestCases = &testcase.StressTest{Runtime: ctrRuntime}
	err := stressTestCases.RunAllTests(context.TODO(), args)
	if err != nil {
		log.Fatal(err)
	}

}
