package steps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	goflow "github.com/kamildrazkiewicz/go-flow"
)

// Step wraps a function and it's dependencies to perform concurrently
// optionally Retry can be specified for steps which need to be retried on error
type Step struct {
	Name      string
	DependsOn []string
	Do        func(context.Context) error
	Retry     int
}

type contextKeyType int

const (
	stepNameKey contextKeyType = iota
)

// Run takes a slice of Steps to be performed concurrently
func Run(ctx context.Context, steps []Step, logger *log.Logger) error {
	flow := goflow.New()
	for _, step := range steps {
		step := step
		dependsOn := step.DependsOn
		if len(dependsOn) == 0 {
			dependsOn = nil
		}
		flow.Add(step.Name, dependsOn, func(r map[string]interface{}) (interface{}, error) {
			if step.Do != nil {
				namedCtx := context.WithValue(ctx, stepNameKey, step.Name)
				l := ContextLogger(namedCtx, logger, "[Steps]")
				attempt := 1
				safelyDo := func() (err error) {
					defer func() {
						if r := recover(); r != nil {
							switch x := r.(type) {
							case string:
								err = errors.New(x)
							case error:
								err = x
							default:
								err = errors.New("Unknown panic")
							}
						}
					}()
					err = step.Do(namedCtx)
					return
				}
				for {
					err := safelyDo()
					if err != nil {
						if attempt <= step.Retry {
							l.Printf("Attempt %d retrying error: %s", attempt, err.Error())
							attempt++
							continue
						}
						l.Printf("Max retry %d reached giving up: %s", attempt, err.Error())
						return nil, fmt.Errorf("step %s failed: %v", step.Name, err)
					}
					return nil, nil
				}
			}
			return nil, nil
		})
	}
	_, err := flow.Do()
	return err
}

// ContextLogger can be used to get a logger prefixed with the current Step Name
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
