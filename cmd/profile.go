package cmd

import (
	"context"
	"log"

	"github.com/kunalkushwaha/ctr-powertest/testcase"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "profile container operations.",

	Run: func(cmd *cobra.Command, args []string) {
		initTestSuite(cmd)

		profileTestCases := &testcase.ProfileContainerTest{Runtime: ctrRuntime}
		err := profileTestCases.RunTestCases(context.TODO(), nil, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(profileCmd)
}
