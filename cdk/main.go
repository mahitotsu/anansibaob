package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdsql"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewAnansibaobStack(scope constructs.Construct, id *string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, id, props)

	dsqlCluster := awsdsql.NewCfnCluster(stack, jsii.String("dsqlCluster"), &awsdsql.CfnClusterProps{
		DeletionProtectionEnabled: false,
	})

	webapi := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("webapi"), &awscdklambdagoalpha.GoFunctionProps{
		Entry: jsii.String("../cmd/webapi"),
		Layers: &[]awslambda.ILayerVersion{
			awslambda.LayerVersion_FromLayerVersionArn(
				stack, jsii.String("WebAdapterLayer"),
				jsii.Sprintf("arn:aws:lambda:%s:753240598075:layer:LambdaAdapterLayerX86:25", *stack.Region()),
			),
		},
		Environment: &map[string]*string{
			"PORT": jsii.String("8000"),
		},
		LogGroup: awslogs.NewLogGroup(stack, jsii.String("LogGroup"), &awslogs.LogGroupProps{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
			Retention:     awslogs.RetentionDays_ONE_DAY,
		}),
	})
	webapi.Role().GrantPrincipal().AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   &[]*string{jsii.String("dsql:DbConnect"), jsii.String("dsql:DbConnectAdmin")},
		Resources: &[]*string{dsqlCluster.AttrResourceArn()},
	}))

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
