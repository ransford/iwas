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
	policiesCmd.Flags().BoolVar(&allPolicies, "all", false, "Show all policies, including AWS managed")
}

var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List IAM policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showPolicies()
	},
}

var allPolicies bool

func showPolicies() error {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := iam.NewFromConfig(cfg)

	// List policies
	listPoliciesInput := &iam.ListPoliciesInput{
		Scope: "Local",
	}
	if allPolicies {
		listPoliciesInput.Scope = "All"
	}
	policies, err := client.ListPolicies(
		context.Background(),
		listPoliciesInput,
	)
	if err != nil {
		return err
	}
	pols := []string{}
	for _, pol := range policies.Policies {
		pols = append(pols, aws.ToString(pol.Arn))
	}
	sort.Strings(pols)
	for _, pol := range pols {
		fmt.Println(pol)
	}

	return nil
}
