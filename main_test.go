package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func argsForTestNoOptional(configFilepath string) []string {
	return []string{
		"./hedwig-terraform-generator",
		"generate",
		configFilepath,
	}
}

func argsForTest(configFilepath string) []string {
	return []string{
		"./hedwig-terraform-generator",
		"generate",
		configFilepath,
		"--alerting",
		"--iam",
		fmt.Sprintf(`--%s=12345`, awsAccountIDFlag),
		fmt.Sprintf(`--%s=us-west-2`, awsRegionFlag),
		fmt.Sprintf(`--%s=pager_action`, queueAlertAlarmActionsFlag),
		fmt.Sprintf(`--%s=pager_action2`, queueAlertAlarmActionsFlag),
		fmt.Sprintf(`--%s=pager_action`, queueAlertOKActionsFlag),
		fmt.Sprintf(`--%s=pager_action2`, queueAlertOKActionsFlag),
		fmt.Sprintf(`--%s=pager_action`, dlqAlertAlarmActionsFlag),
		fmt.Sprintf(`--%s=pager_action2`, dlqAlertAlarmActionsFlag),
		fmt.Sprintf(`--%s=pager_action`, dlqAlertOKActionsFlag),
		fmt.Sprintf(`--%s=pager_action2`, dlqAlertOKActionsFlag),
	}
}

func TestGenerate(t *testing.T) {
	info, err := ioutil.ReadDir("test_fixtures")
	require.NoError(t, err)

	dmp := diffmatchpatch.New()

	for _, testDir := range info {
		if !testDir.IsDir() {
			continue
		}

		os.RemoveAll("hedwig")

		fmt.Println("Test:", testDir.Name())

		configFilepath := filepath.Join("test_fixtures", testDir.Name(), "test_config.json")

		var args []string
		if strings.Contains(testDir.Name(), "no_optional_param") {
			args = argsForTestNoOptional(configFilepath)
		} else {
			args = argsForTest(configFilepath)
		}

		assert.NoError(t, runApp(args))

		info, err := ioutil.ReadDir("hedwig")
		assert.NoError(t, err)

		files := make([]string, len(info))
		for i, f := range info {
			files[i] = f.Name()
		}

		infoTestDir, err := ioutil.ReadDir(filepath.Join("test_fixtures", testDir.Name()))
		require.NoError(t, err)

		var testFiles []string
		for _, testOutputFile := range infoTestDir {
			if filepath.Ext(testOutputFile.Name()) != ".tf" {
				continue
			}
			testFiles = append(testFiles, testOutputFile.Name())
		}

		assert.Equal(t, testFiles, files)

		for _, testOutputFile := range infoTestDir {
			if filepath.Ext(testOutputFile.Name()) != ".tf" {
				continue
			}
			testOutputFileName := filepath.Join("test_fixtures", testDir.Name(), testOutputFile.Name())
			expectedBytes, err := ioutil.ReadFile(testOutputFileName)
			require.NoError(t, err)

			expected := string(expectedBytes)

			// poor template engine
			expected = strings.Replace(
				expected, "{{TFQueueModuleVersion}}", TFQueueModuleVersion, -1)
			expected = strings.Replace(
				expected, "{{TFQueueSubscriptionModuleVersion}}", TFQueueSubscriptionModuleVersion, -1)
			expected = strings.Replace(
				expected, "{{TFLambdaSubscriptionModuleVersion}}", TFLambdaSubscriptionModuleVersion, -1)
			expected = strings.Replace(
				expected, "{{TFTopicModuleVersion}}", TFTopicModuleVersion, -1)
			expected = strings.Replace(
				expected, "{{GENERATOR_VERSION}}", VERSION, -1)

			actualB, err := ioutil.ReadFile(filepath.Join("hedwig", testOutputFile.Name()))
			require.NoError(t, err)

			assert.Equal(
				t, expected, string(actualB),
				dmp.DiffPrettyText(dmp.DiffMain(expected, string(actualB), true)),
			)
		}

		if t.Failed() {
			// so we can investigate what went wrong
			break
		}
	}
}
