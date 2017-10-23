package cmd

import (
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

func init() {
	RootCmd.AddCommand(basicCmd)

	//Set Log level
	log.SetLevel(log.DebugLevel)
}

func runBasicTest(cmd *cobra.Command, args []string) {
	initTestSuite(cmd)

	singleClientTestCases := &testcase.BasicContainerTest{Runtime: ctrRuntime}

	err := singleClientTestCases.RunTestCases(ctx, nil, args)
	if err != nil {
		log.Fatal(err)
	}
}
