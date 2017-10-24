package cmd

import (
	"github.com/kunalkushwaha/ctr-powertest/testcase"

	"log"

	"github.com/spf13/cobra"
)

// parallelCmd represents the parallel command
var stressCmd = &cobra.Command{
	Use:     "stress",
	Short:   "Run container tests in parallel (Stress Test)",
	Example: "sudo ctr-powertest -p containerd stress image-pull",
	Run:     runStressTest,
}

func init() {
	RootCmd.AddCommand(stressCmd)

	stressCmd.Flags().StringSliceP("testcase", "t", nil, "Testcases to run [container-create-delete | image-pull]")
}

func runStressTest(cmd *cobra.Command, args []string) {
	initTestSuite(cmd)

	testcases, _ := cmd.Flags().GetStringSlice("testcase")
	stressTestCases := &testcase.StressTest{Runtime: ctrRuntime}
	err := stressTestCases.RunTestCases(ctx, testcases, args)
	if err != nil {
		log.Fatal(err)
	}

}
