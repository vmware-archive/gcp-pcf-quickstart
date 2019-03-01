package parser

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func NewUint(set *flag.FlagSet, field reflect.Value, tags reflect.StructTag) (*Flag, error) {
	var defaultValue uint64
	defaultStr, ok := tags.Lookup("default")
	if ok {
		var err error
		defaultValue, err = strconv.ParseUint(defaultStr, 0, 0)
		if err != nil {
			return &Flag{}, fmt.Errorf("could not parse uint default value %q: %s", defaultStr, err)
		}
	}

	var f Flag
	short, ok := tags.Lookup("short")
	if ok {
		set.UintVar(field.Addr().Interface().(*uint), short, uint(defaultValue), "")
		f.flags = append(f.flags, set.Lookup(short))
		f.name = fmt.Sprintf("-%s", short)
	}

	long, ok := tags.Lookup("long")
	if ok {
		set.UintVar(field.Addr().Interface().(*uint), long, uint(defaultValue), "")
		f.flags = append(f.flags, set.Lookup(long))
		f.name = fmt.Sprintf("--%s", long)
	}

	alias, ok := tags.Lookup("alias")
	if ok {
		set.UintVar(field.Addr().Interface().(*uint), alias, uint(defaultValue), "")
		f.flags = append(f.flags, set.Lookup(alias))
		f.name = fmt.Sprintf("--%s", alias)
	}

	env, ok := tags.Lookup("env")
	if ok {
		envOpts := strings.Split(env, ",")

		for _, envOpt := range envOpts {
			envStr := os.Getenv(envOpt)
			if envStr != "" {
				envValue, err := strconv.ParseUint(envStr, 0, 0)
				if err != nil {
					return &Flag{}, fmt.Errorf("could not parse uint environment variable %s value %q: %s", envOpt, envStr, err)
				}

				field.SetUint(envValue)
				f.set = true
				break
			}
		}
	}

	_, f.required = tags.Lookup("required")

	return &f, nil
}
