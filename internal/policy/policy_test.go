package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const st1 = `{
	"Version": "2012-10-17",
	"Statement": [
	  {
		"Sid": "DenyAllAwsResourcesOutsideAccountExceptSNS",
		"Effect": "Deny",
		"NotAction": [
		  "sns:*"
		],
		"Resource": "*",
		"Condition": {
		  "StringNotEquals": {
			"aws:ResourceAccount": [
			  "111122223333"
			]
		  }
		}
	  },
	  {
		"Sid": "DenyAllSNSResourcesOutsideAccountExceptCloudFormation",
		"Effect": "Deny",
		"Action": [
		  "sns:*"
		],
		"Resource": "*",
		"Condition": {
		  "StringNotEquals": {
			"aws:ResourceAccount": [
			  "111122223333"
			]
		  },
		  "ForAllValues:StringNotEquals": {
			"aws:CalledVia": [
			  "cloudformation.amazonaws.com"
			]
		  }
		}
	  }
	]
  }`

func TestNewFromJSON(t *testing.T) {
	p, err := NewFromJson([]byte(st1))
	assert.NoError(t, err)
	assert.Equal(t, p.Version, "2012-10-17")
	assert.Len(t, p.Statements, 2)
}
