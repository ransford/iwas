package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(policiesCmd)
}

var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List IAM policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showPolicies()
	},
}

func showPolicies() error {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := iam.NewFromConfig(cfg)

	// List policies
	policies, err := client.ListPolicies(
		context.Background(),
		&iam.ListPoliciesInput{
			Scope: "Local",
		},
	)
	if err != nil {
		return err
	}
	pols := []string{}
	for _, pol := range policies.Policies {
		pols = append(pols, aws.ToString(pol.PolicyName))
	}
	sort.Strings(pols)
	for _, pol := range pols {
		fmt.Println(pol)
	}

	return nil
}
