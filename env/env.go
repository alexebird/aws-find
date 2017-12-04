package env

import (
	"fmt"
	"os"
)

func GetEnv(envVarName string) *string {
	env := os.Getenv(envVarName)
	if env == "" {
		return nil
	} else {
		return &env
	}
}

func MustGetEnv(envVarName string) *string {
	env := os.Getenv(envVarName)
	if env == "" {
		fmt.Println("must set", envVarName)
		os.Exit(1)
		return nil
	} else {
		return &env
	}
}

func DavinciEnv() *string {
	return GetEnv("DAVINCI_ENV")
}

func DavinciEnvFull() *string {
	return GetEnv("DAVINCI_ENV_FULL")
}
