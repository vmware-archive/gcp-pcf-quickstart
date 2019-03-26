package parser

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func NewDuration(set *flag.FlagSet, field reflect.Value, tags reflect.StructTag) (*Flag, error) {
	var defaultValue time.Duration
	defaultStr, ok := tags.Lookup("default")
	if ok {
		var err error
		defaultValue, err = time.ParseDuration(defaultStr)
		if err != nil {
			return &Flag{}, fmt.Errorf("could not parse duration default value %q: %s", defaultStr, err)
		}
	}

	var f Flag
	short, ok := tags.Lookup("short")
	if ok {
		set.DurationVar(field.Addr().Interface().(*time.Duration), short, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(short))
		f.name = fmt.Sprintf("-%s", short)
	}

	long, ok := tags.Lookup("long")
	if ok {
		set.DurationVar(field.Addr().Interface().(*time.Duration), long, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(long))
		f.name = fmt.Sprintf("--%s", long)
	}

	alias, ok := tags.Lookup("alias")
	if ok {
		set.DurationVar(field.Addr().Interface().(*time.Duration), alias, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(alias))
		f.name = fmt.Sprintf("--%s", alias)
	}

	env, ok := tags.Lookup("env")
	if ok {
		envOpts := strings.Split(env, ",")

		for _, envOpt := range envOpts {
			envStr := os.Getenv(envOpt)
			if envStr != "" {
				envValue, err := time.ParseDuration(envStr)
				if err != nil {
					return &Flag{}, fmt.Errorf("could not parse duration environment variable %s value %q: %s", envOpt, envStr, err)
				}

				field.SetInt(int64(envValue))
				f.set = true
				break
			}
		}
	}

	_, f.required = tags.Lookup("required")

	return &f, nil
}
