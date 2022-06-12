package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	pols := []string{}
	for _, pol := range policies.Policies {
		pols = append(pols, aws.ToString(pol.PolicyName))
	}
	sort.Strings(pols)
	for _, pol := range pols {
		fmt.Println(pol)
	}
}
