package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"os"
)

func TestRetrieveExportableEnvironmentSingleLength(t *testing.T) {
	environments := createExportedEnvironment("HOME")
	require.Equal(t, len(environments), 1)
}

func TestRetrieveExportableEnvironmentMultipleLength(t *testing.T) {
	environments := createExportedEnvironment("HOME|USER")
	require.Equal(t, len(environments), 2)
}

func TestRetrieveExportableEnvironmentValues(t *testing.T) {
	environments := createExportedEnvironment("HOME|USER")
	require.Equal(t, environments[0].MappedTo, "HOME")
	require.Equal(t, environments[0].Value, os.Getenv("HOME"))
	require.Equal(t, environments[1].MappedTo, "USER")
	require.Equal(t, environments[1].Value, os.Getenv("USER"))
}

func TestRetrieveExportableEnvironmentEmptyValue(t *testing.T) {
	environments := createExportedEnvironment("dummy")
	require.Equal(t, environments[0].MappedTo, "dummy")
	require.Equal(t, environments[0].Value, "")
}

func TestSplitPipeSeparatedStringArrayEmpty(t *testing.T) {
	environment := splitPipeSeparatedStringArray("")
	require.Equal(t, 0, len(environment))
}

func TestSplitPipeSeparatedStringArrayNonEmpty(t *testing.T) {
	environment := splitPipeSeparatedStringArray("a|b")
	require.Equal(t, 2, len(environment))
	require.Equal(t, environment[0], "a")
	require.Equal(t, environment[1], "b")
}

func TestValidateConfigsNoSlug(t *testing.T) {
	configs := ConfigsModel{
		APIToken:"token",
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsNoApiToken(t *testing.T) {
	configs := ConfigsModel{
		AppSlug:"slug",
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsEmptyVariable(t *testing.T) {
	configs := ConfigsModel{
		APIToken:"token",
		AppSlug:"slug",
		ExportedVariableNames:"|",
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsInvalidVariable(t *testing.T) {
	configs := ConfigsModel{
		APIToken:"token",
		AppSlug:"slug",
		ExportedVariableNames:"a|b=",
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsNoVariables(t *testing.T) {
	configs := ConfigsModel{
		APIToken:"token",
		AppSlug:"slug",
		ExportedVariableNames:"",
	}
	require.NoError(t, configs.validate())
}

func TestValidateConfigsValidVariables(t *testing.T) {
	configs := ConfigsModel{
		APIToken:"token",
		AppSlug:"slug",
		ExportedVariableNames:"a|b",
	}
	require.NoError(t, configs.validate())
}
