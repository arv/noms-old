// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package dynamodb

import (
	"github.com/attic-labs/noms/Godeps/_workspace/src/github.com/aws/aws-sdk-go/private/waiter"
)

func (c *DynamoDB) WaitUntilTableExists(input *DescribeTableInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeTable",
		Delay:       20,
		MaxAttempts: 25,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "path",
				Argument: "Table.TableStatus",
				Expected: "ACTIVE",
			},
			{
				State:    "retry",
				Matcher:  "error",
				Argument: "",
				Expected: "ResourceNotFoundException",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

func (c *DynamoDB) WaitUntilTableNotExists(input *DescribeTableInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeTable",
		Delay:       20,
		MaxAttempts: 25,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "error",
				Argument: "",
				Expected: "ResourceNotFoundException",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}
