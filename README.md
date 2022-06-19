# iwas

Versioned policy browser for AWS IAM policies.

# Usage

Show all policies:

    iwas policies
    iwas list

Show the current version of a policy:

    iwas policy <policy ARN>
    iwas get <policy ARN>

Show the history of a policy:

    iwas history <policy ARN>
    iwas log <policy ARN>

Show a previous version of a policy:

    iwas policy <policy ARN> <version>
    iwas get <policy ARN> <version>
