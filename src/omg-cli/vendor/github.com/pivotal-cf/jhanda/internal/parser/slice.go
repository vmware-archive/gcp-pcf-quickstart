package parser

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func NewSlice(set *flag.FlagSet, field reflect.Value, tags reflect.StructTag) (*Flag, error) {
	collection := field.Addr().Interface().(*[]string)

	defaultSlice, ok := tags.Lookup("default")
	if ok {
		separated := strings.Split(defaultSlice, ",")
		*collection = append(*collection, separated...)
	}

	slice := StringSlice{collection}

	var f Flag
	short, ok := tags.Lookup("short")
	if ok {
		set.Var(&slice, short, "")
		f.flags = append(f.flags, set.Lookup(short))
		f.name = fmt.Sprintf("-%s", short)
	}

	long, ok := tags.Lookup("long")
	if ok {
		set.Var(&slice, long, "")
		f.flags = append(f.flags, set.Lookup(long))
		f.name = fmt.Sprintf("--%s", long)
	}

	alias, ok := tags.Lookup("alias")
	if ok {
		set.Var(&slice, alias, "")
		f.flags = append(f.flags, set.Lookup(alias))
		f.name = fmt.Sprintf("--%s", alias)
	}

	env, ok := tags.Lookup("env")
	if ok {
		envOpts := strings.Split(env, ",")

		for _, envOpt := range envOpts {
			envStr, ok := os.LookupEnv(envOpt)
			if ok {
				separated := strings.Split(envStr, ",")
				*collection = append(*collection, separated...)
				f.set = true
				break
			}
		}
	}

	_, f.required = tags.Lookup("required")

	return &f, nil
}

type StringSlice struct {
	slice *[]string
}

func (ss *StringSlice) String() string {
	return fmt.Sprintf("%s", ss.slice)
}

func (ss *StringSlice) Set(item string) error {
	*ss.slice = append(*ss.slice, item)
	return nil
}
