package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.Flags().BoolVar(&allPolicies, "all", false, "Show all policies, including AWS managed")
}

var policiesCmd = &cobra.Command{
	Use:     "policies",
	Aliases: []string{"list"},
	Short:   "List IAM policies",
	PreRun:  setLogLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showPolicies()
	},
}

var allPolicies bool

func showPolicies() error {
	// List policies
	listPoliciesInput := &iam.ListPoliciesInput{
		Scope: "Local",
	}
	if allPolicies {
		listPoliciesInput.Scope = "All"
	}
	log.Debug("Listing policies")

	paginator := iam.NewListPoliciesPaginator(iamClient, listPoliciesInput, func(o *iam.ListPoliciesPaginatorOptions) {
		o.Limit = 100
	})

	numPolicies := 0
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return err
		}
		for _, pol := range output.Policies {
			fmt.Println(aws.ToString(pol.Arn))
			numPolicies += 1
		}
	}
	log.Debugf("Listed %d policies", numPolicies)

	return nil
}
