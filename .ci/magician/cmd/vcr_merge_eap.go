package cmd

import (
	"fmt"
	"magician/exec"
	"magician/source"
	"os"

	"github.com/spf13/cobra"
)

var vcrMergeEapCmd = &cobra.Command{
	Use:   "vcr-merge-eap",
	Short: "Merge VCR cassettes for EAP",
	Long: `This command is triggered in .ci/gcb-push-downstream.yml to merge vcr cassettes.

	The command expects the following as arguments:
	1. CL number

	It then performs the following operations:
	1. Run gsutil to list, copy, and remove the vcr cassettes fixtures.
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clNumber := args[0]
		fmt.Println("CL number:", clNumber)

		baseBranch := os.Getenv("BASE_BRANCH")
		if baseBranch == "" {
			return fmt.Errorf("environment variable BASE_BRANCH is empty")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %w", err)
		}

		return execVCRMergeEAP(clNumber, baseBranch, rnr)
	},
}

func execVCRMergeEAP(clNumber, baseBranch string, runner source.Runner) error {
	head := "auto-cl-" + clNumber
	mergeCassettes("gs://ci-vcr-cassettes/private", baseBranch, fmt.Sprintf("refs/heads/%s", head), runner)
	return nil
}

func init() {
	rootCmd.AddCommand(vcrMergeEapCmd)
}
