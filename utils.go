package main

import (
	"github.com/bitrise-io/go-utils/command"
	"strings"
	"os"
)

func exportEnvironmentWithEnvman(keyStr, valueStr string) error {
	cmd := command.New("envman", "add", "--key", keyStr)
	cmd.SetStdin(strings.NewReader(valueStr))
	return cmd.Run()
}

func createExportedEnvironment(variableNames string) []EnvironmentVariableModel {
	environmentVariableNames := splitPipeSeparatedStringArray(variableNames)
	environments := []EnvironmentVariableModel{}

	for _, environmentVariableName := range environmentVariableNames {
		environments = append(environments, EnvironmentVariableModel{
			MappedTo:"",
			Value:os.Getenv(environmentVariableName),
		})
	}
	return environments
}

func splitPipeSeparatedStringArray(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, "|")
}