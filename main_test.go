package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"os"
)

func TestRetrieveExportableEnvironmentSingleLength(t *testing.T) {
	environments := createExportedEnvironment("HOME")
	require.Equal(t, 1, len(environments))
}

func TestRetrieveExportableEnvironmentMultipleLength(t *testing.T) {
	environments := createExportedEnvironment("HOME|USER")
	require.Equal(t, 2, len(environments))
}

func TestRetrieveExportableEnvironmentValues(t *testing.T) {
	environments := createExportedEnvironment("HOME|USER")
	require.Equal(t, "HOME", environments[0].MappedTo)
	require.Equal(t, os.Getenv("HOME"), environments[0].Value)
	require.Equal(t, "USER", environments[1].MappedTo)
	require.Equal(t, os.Getenv("USER"), environments[1].Value)
}

func TestRetrieveExportableEnvironmentEmptyValue(t *testing.T) {
	environments := createExportedEnvironment("dummy")
	require.Equal(t, "dummy", environments[0].MappedTo)
	require.Equal(t, "", environments[0].Value)
}

func TestSplitPipeSeparatedStringArrayEmpty(t *testing.T) {
	environment := splitPipeSeparatedStringArray("")
	require.Equal(t, 0, len(environment))
}

func TestSplitPipeSeparatedStringArrayNonEmpty(t *testing.T) {
	environment := splitPipeSeparatedStringArray("a|b")
	require.Equal(t, 2, len(environment))
	require.Equal(t, "a", environment[0])
	require.Equal(t, "b", environment[1])
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
