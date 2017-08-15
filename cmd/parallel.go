package cmd

import (
	//	"context"
	//	"log"

	"context"
	"log"

	"github.com/kunalkushwaha/ctr-powertest/testcase"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parallelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parallelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func runParallelTest(cmd *cobra.Command, args []string) {

	testSetup, err := testcase.SetupParallelTestEnvironment(config.RuntimeName, config, false)
	if err != nil {
		log.Fatal("Error while setting up environment : ", err)
	}

	err = testSetup.RunAllTests(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
