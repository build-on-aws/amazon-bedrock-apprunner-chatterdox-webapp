package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdkapprunneralpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const appPort = 8080

type ChatterdoxStackProps struct {
	awscdk.StackProps
}

func main() {
	app := awscdk.NewApp(nil)
	ChatterdoxStack(app, "ChatterdoxStack", &ChatterdoxStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

const appDir = "../chatterdox-web-app"

func ChatterdoxStack(scope constructs.Construct, id string, props *ChatterdoxStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	appImage := awsecrassets.NewDockerImageAsset(stack, jsii.String("app-image"),
		&awsecrassets.DockerImageAssetProps{
			Directory: jsii.String(appDir)})

	apprunnerInstanceRole := awsiam.NewRole(stack, jsii.String("apprunner-bedrock-role"),
		&awsiam.RoleProps{
			AssumedBy: awsiam.NewServicePrincipal(jsii.String("tasks.apprunner.amazonaws.com"), nil),
		})

	apprunnerInstanceRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("bedrock:*"),
		Resources: jsii.Strings("*")}))

	app := awscdkapprunneralpha.NewService(stack, jsii.String("chatterdox_app"),
		&awscdkapprunneralpha.ServiceProps{
			ServiceName: jsii.String("chatterdox"),
			Source: awscdkapprunneralpha.NewAssetSource(
				&awscdkapprunneralpha.AssetProps{
					ImageConfiguration: &awscdkapprunneralpha.ImageConfiguration{
						Port: jsii.Number(appPort)},
					Asset: appImage}),
			InstanceRole: apprunnerInstanceRole,
			Memory:       awscdkapprunneralpha.Memory_TWO_GB(),
			Cpu:          awscdkapprunneralpha.Cpu_ONE_VCPU(),
		})

	app.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	appurl := "https://" + *app.ServiceUrl()
	awscdk.NewCfnOutput(stack, jsii.String("AppURL"), &awscdk.CfnOutputProps{Value: jsii.String(appurl)})

	return stack
}

func env() *awscdk.Environment {
	return nil

}
