package utils

import (
	"flag"
	"os"
	"strconv"
)

func GetStringParam(envName string, flagName string, flagUsage string, def string) *string {
	value := os.Getenv(envName)
	if value != "" {
		return &value
	}
	return flag.String(flagName, def, flagUsage)
}

func GetIntParam(envName string, flagName string, flagUsage string, def int) *int {
	value := os.Getenv(envName)
	if value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return &i
	}
	return flag.Int(flagName, def, flagUsage)
}
