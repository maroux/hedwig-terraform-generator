package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

const (
	// VERSION represents the version of the generator tool
	VERSION = "1.2.0"

	// TFQueueModuleVersion represents the version of the hedwig-queue module
	TFQueueModuleVersion = "1.0.0"

	// TFQueueSubscriptionModuleVersion represents the version of the hedwig-queue-subscription module
	TFQueueSubscriptionModuleVersion = "1.0.0"

	// TFLambdaSubscriptionModuleVersion represents the version of the hedwig-lambda-subscription module
	TFLambdaSubscriptionModuleVersion = "1.0.0"

	// TFTopicModuleVersion represents the version of the hedwig-topic module
	TFTopicModuleVersion = "1.0.0"

	tfDoNotEditStamp = `// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator ` + VERSION
)

const (
	// alertingFlag represents the cli flag that indicates if alerting should be generated
	alertingFlag = "alerting"

	// awsAccountIDFlag represents the cli flag for aws account id
	awsAccountIDFlag = "aws-account-id"

	// awsRegionFlag represents the cli flag for aws region
	awsRegionFlag = "aws-region"

	// dlqAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	dlqAlertAlarmActionsFlag = "dlq-alert-alarm-actions"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	dlqAlertOKActionsFlag = "dlq-alert-ok-actions"

	// moduleFlag represents the cli flag for output module name
	moduleFlag = "module"

	// queueAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	queueAlertAlarmActionsFlag = "queue-alert-alarm-actions"

	// queueAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	queueAlertOKActionsFlag = "queue-alert-ok-actions"
)

func validateArgs(c *cli.Context) *cli.ExitError {
	if c.NArg() != 1 {
		return cli.NewExitError("<config-file> is required", 1)
	}

	alertingFlagsOkay := true
	alertingFlags := []string{queueAlertAlarmActionsFlag, queueAlertOKActionsFlag, dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag}
	if c.Bool(alertingFlag) {
		for _, f := range alertingFlags {
			if !c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is required", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("missing required flags for --alerting", 1)
		}
	} else {
		for _, f := range alertingFlags {
			if c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is disallowed", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("disallowed flags specified with missing --alerting", 1)
		}
	}
	return nil
}

func generateModule(c *cli.Context) error {
	if err := validateArgs(c); err != nil {
		return err
	}

	configFile := c.Args().Get(0)

	config, err := NewConfig(configFile)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = writeTerraform(config, c)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "failed to generate terraform module"), 1)
	}

	fmt.Println("Created Terraform Hedwig module successfully!")
	return nil
}

func generateConfigFileStructure(c *cli.Context) error {
	structure := Config{
		Topics: []string{
			"my-topic",
		},
		QueueConsumers: []*QueueConsumer{
			{
				"DEV-MYAPP",
				map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
				[]string{"my-topic"},
			},
		},
		LambdaConsumers: []*LambdaConsumer{
			{
				"arn:aws:lambda:us-west-2:12345:function:my_function:deployed",
				"{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
				"{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
				[]string{"my-topic"},
			},
		},
	}
	structureAsJSON, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(string(structureAsJSON))
	return nil
}

func runApp(args []string) error {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app := cli.NewApp()
	app.Name = "Hedwig Terraform"
	app.Usage = "Manage Terraform configuration for Hedwig apps"
	app.Version = VERSION
	app.HelpName = "hedwig-terraform"
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			Usage:     "Generates Terraform module for Hedwig apps",
			ArgsUsage: "<config-file>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  moduleFlag,
					Usage: "Terraform module name to generate",
					Value: "hedwig",
				},
				cli.BoolFlag{
					Name:  alertingFlag,
					Usage: "Should Cloudwatch alerting be generated?",
				},
				cli.StringSliceFlag{
					Name:  queueAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when in ALARM",
				},
				cli.StringSliceFlag{
					Name:  queueAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when OK",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when in ALARM",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when OK",
				},
				cli.StringFlag{
					Name:  awsAccountIDFlag,
					Usage: "AWS Account ID",
				},
				cli.StringFlag{
					Name:  awsRegionFlag,
					Usage: "AWS Region",
				},
			},
			Action: generateModule,
		},
		{
			Name:   "config-file-structure",
			Usage:  "Outputs the structure for config file required for generate command",
			Action: generateConfigFileStructure,
		},
	}

	return app.Run(args)
}

func main() {
	if err := runApp(os.Args); err != nil {
		log.Fatal(err)
	}
}
