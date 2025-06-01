package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewAnansibaobStack(scope constructs.Construct, id *string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, id, props)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAnansibaobStack(app, jsii.String("AnansibaobStack"), &awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			Region:  jsii.String("ap-northeast-1"),
		},
	})

	app.Synth(nil)
}
