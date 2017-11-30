package davinci

import (
	"os"
)

func getEnv(envVarName string) *string {
	env := os.Getenv(envVarName)
	if env == "" {
		return nil
	} else {
		return &env
	}
}

func DavinciEnv() *string {
	return getEnv("DAVINCI_ENV")
}

func DavinciEnvFull() *string {
	return getEnv("DAVINCI_ENV_FULL")
}
