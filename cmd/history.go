package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ransford/iwas/internal"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

func getPolicyHistory(arn string) error {
	getPolicyInput := iam.GetPolicyInput{
		PolicyArn: &arn,
	}
	pol, err := iamClient.GetPolicy(context.Background(), &getPolicyInput)
	if err != nil {
		return err
	}
	ver, err := internal.ParseVersion(*pol.Policy.DefaultVersionId)
	if err != nil {
		return err
	}
	for i := ver; i >= 1; i-- {
		doc, err := getPolicyVersion(arn, fmt.Sprintf("v%d", i))
		if err != nil {
			return err
		}
		fmt.Println(doc)
	}
	return nil
}

var historyCmd = &cobra.Command{
	Use:     "history <arn>",
	Aliases: []string{"get"},
	Short:   "Show all versions of an IAM policy",
	Args: func(cmd *cobra.Command, args []string) error {
		// ARN only
		if len(args) == 1 {
			if _, err := internal.PolicyNameToArn(args[0]); err != nil {
				return err
			}
			return nil
		}
		return errors.New("wrong number of arguments")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		parn, _ := internal.PolicyNameToArn(args[0])
		log.Println(parn.String())
		return getPolicyHistory(parn.String())
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
