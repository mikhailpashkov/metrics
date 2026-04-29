package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ValueConsumer[T string | int | bool] func(T)

type PendingFlagParse func()

type Param interface {
	do() PendingFlagParse
}

type StringParam struct {
	EnvName       string
	FlagName      string
	FlagUsage     string
	Default       string
	ValueConsumer ValueConsumer[string]
}

type IntParam struct {
	EnvName       string
	FlagName      string
	FlagUsage     string
	Default       int
	ValueConsumer ValueConsumer[int]
}

type BoolParam struct {
	EnvName       string
	FlagName      string
	FlagUsage     string
	Default       bool
	ValueConsumer ValueConsumer[bool]
}

func (param *StringParam) do() PendingFlagParse {
	value := os.Getenv(param.EnvName)
	if value != "" {
		param.ValueConsumer(value)
		return nil
	}
	valuePtr := flag.String(param.FlagName, param.Default, param.FlagUsage)
	return func() {
		param.ValueConsumer(*valuePtr)
	}
}

func (param *IntParam) do() PendingFlagParse {
	value := os.Getenv(param.EnvName)
	if value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		param.ValueConsumer(i)
		return nil
	}
	valuePtr := flag.Int(param.FlagName, param.Default, param.FlagUsage)
	return func() {
		param.ValueConsumer(*valuePtr)
	}
}

func (param *BoolParam) do() PendingFlagParse {
	value := os.Getenv(param.EnvName)
	if value != "" {
		switch value {
		case "true":
			param.ValueConsumer(true)
		case "false":
			param.ValueConsumer(false)
		default:
			panic(fmt.Errorf("invalid boolean value for env %s: %s", param.EnvName, value))
		}
		return nil
	}
	valuePtr := flag.Bool(param.FlagName, param.Default, param.FlagUsage)
	return func() {
		param.ValueConsumer(*valuePtr)
	}
}

func GetParams(configs []Param) {
	pendingFlagParseList := make([]PendingFlagParse, 0, len(configs))
	for _, config := range configs {
		pendingFlagParseList = append(pendingFlagParseList, config.do())
	}
	flag.Parse()
	for _, pendingFlagParse := range pendingFlagParseList {
		if pendingFlagParse == nil {
			continue
		}
		pendingFlagParse()
	}
}
