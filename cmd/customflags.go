package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type AllowedStringValuesFlag struct {
	*cli.StringFlag
	AllowedValues []string
	AllowEmpty    bool
}

func (f *AllowedStringValuesFlag) Apply(set *flag.FlagSet) error {
	if err := f.StringFlag.Apply(set); err != nil {
		return err
	}

	if allowed, err := checkAllowedValue(f.Name, f.StringFlag.Value, f.AllowedValues); !allowed && err != nil {
		if !(f.StringFlag.Value == "" && f.AllowEmpty) {
			return err
		}
	}

	return nil
}

func checkAllowedValue[T comparable](flagName string, flagVal T, allowedValues []T) (allowed bool, err error) {
	for _, val := range allowedValues {
		if flagVal == val {
			allowed = true
			break
		}
	}

	if !allowed {
		err = errors.Errorf(
			"value %[1]v not allowed for flag %[2]s.\n"+
				"allowed values: (%[3]s)",
			flagVal, flagName, allowedValuesSlice(allowedValues),
		)
	}

	return
}

func allowedValuesSlice[T any](allowedValues []T) (out string) {
	strVals := make([]string, len(allowedValues))

	for i, val := range allowedValues {
		switch v := (any)(val).(type) {
		case int:
			strVals[i] = strconv.Itoa(v)
		case string:
			strVals[i] = v
		}
	}

	out = strings.Join(strVals, " ")

	return
}
