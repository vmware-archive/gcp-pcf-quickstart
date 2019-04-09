package steps

import (
	"context"
	"log"
	"strings"
)

type contextKeyType int

const (
	stepNameKey contextKeyType = iota
)

func Step(ctx context.Context, name string, f func(context.Context) error) func(map[string]interface{}) (interface{}, error) {
	return func(r map[string]interface{}) (interface{}, error) {
		namedCtx := context.WithValue(ctx, stepNameKey, name)
		err := f(namedCtx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func ContextLogger(ctx context.Context, logger *log.Logger, prefix string) *log.Logger {
	prefixes := []string{}
	if logger.Prefix() != "" {
		prefixes = append(prefixes, logger.Prefix())
	}
	if prefix != "" {
		prefixes = append(prefixes, prefix)
	}
	if v, ok := ctx.Value(stepNameKey).(string); ok {
		prefixes = append(prefixes, v)
	}
	if len(prefixes) > 0 {
		prefixes = append(prefixes, "")
	}

	return log.New(logger.Writer(), strings.Join(prefixes, " "), 0)
}
