package util

import (
	"os"
	"strconv"
)

func LookupEnvAsBool(varName string) bool {
	varValue, ok := os.LookupEnv(varName)
	if !ok {
		return false
	}
	b, err := strconv.ParseBool(varValue)
	if err != nil {
		return false
	}
	return b
}
