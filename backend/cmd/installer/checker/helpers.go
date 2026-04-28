package checker

import (
	"os"

	"pentagi/cmd/installer/state"
)

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkFileIsReadable(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// checkDirIsWritable checks if we can write to a directory
func checkDirIsWritable(dirPath string) bool {
	// try to create a temporary file in the directory
	tempFile, err := os.CreateTemp(dirPath, ".pentagi_test_*")
	if err != nil {
		return false
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// clean up the test file
	os.Remove(tempPath)
	return true
}

func getEnvVar(appState state.State, key, defaultValue string) string {
	if appState == nil {
		return defaultValue
	}

	if envVar, exist := appState.GetVar(key); exist && envVar.Value != "" {
		return envVar.Value
	} else if envVar.Default != "" {
		return envVar.Default
	}

	return defaultValue
}

// getProxyURL retrieves the proxy URL from application state if configured
func getProxyURL(appState state.State) string {
	if appState == nil {
		return ""
	}
	return getEnvVar(appState, "PROXY_URL", "")
}
