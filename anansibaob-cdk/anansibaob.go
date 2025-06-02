package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewAnansibaobStack(scope constructs.Construct, id *string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, id, props)

	webapi := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("webapi"), &awscdklambdagoalpha.GoFunctionProps{
		Entry: jsii.String("../anansibaob-lambda/webapi"),
		LogGroup: awslogs.NewLogGroup(stack, jsii.String("LogGroup"), &awslogs.LogGroupProps{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
			Retention:     awslogs.RetentionDays_ONE_DAY,
		}),
	})
	weburi := webapi.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		InvokeMode: awslambda.InvokeMode_BUFFERED,
		AuthType:   awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String("endpoint"), &awscdk.CfnOutputProps{
		Value: weburi.Url(),
	})

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
