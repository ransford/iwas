package internal

import "encoding/json"

const (
	Allow PolicyEffect = "Allow"
	Deny  PolicyEffect = "Deny"
)

type Principal map[string]string
type Action string
type Resource string
type PolicyEffect string
type Condition map[string]ConditionElement
type ConditionElement map[string][]string

type Policy struct {
	Version    string      `json:"Version"`
	Id         string      `json:"Id"`
	Statements []Statement `json:"Statement"`
}

// Statement represents an IAM policy statement.
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html
type Statement struct {
	Effect PolicyEffect `json:"Effect"`
	Sid    string       `json:"Sid"`

	Principals    []Principal `json:"Principal"`
	NotPrincipals []Principal `json:"NotPrincipal"`

	Actions    []Action `json:"Action"`
	NotActions []Action `json:"NotAction"`

	Resources    Resource `json:"Resource"`
	NotResources Resource `json:"NotResource"`

	Conditions Condition `json:"Condition"`
}

func NewFromJson(j []byte) (*Policy, error) {
	pol := &Policy{}
	if err := json.Unmarshal(j, pol); err != nil {
		return nil, err
	}
	return pol, nil
}
