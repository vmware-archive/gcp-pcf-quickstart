package steps

import (
	"context"
	"log"
	"strings"

	goflow "github.com/kamildrazkiewicz/go-flow"
)

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
				for {
					err := step.Do(namedCtx)
					if err != nil {
						if attempt <= step.Retry {
							l.Printf("Attempt %d retrying error: %s", attempt, err.Error())
							attempt++
							continue
						}
						l.Printf("Max retry %d reached giving up: %s", attempt, err.Error())
						return nil, err
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
