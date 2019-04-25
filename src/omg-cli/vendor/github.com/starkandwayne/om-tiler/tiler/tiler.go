package tiler

import (
	"context"
	"fmt"
	"log"

	"github.com/starkandwayne/om-tiler/steps"
)

type callbackType string

const (
	// BuildCallback register Step to be invoked during Build
	BuildCallback callbackType = "BuildCallback"
	// DeleteCallback register Step to be invoked during Delete
	DeleteCallback = "DeleteCallback"
)

var allowedCallbackHooks = map[callbackType][]string{
	BuildCallback: {
		StepWaitOpsmanOnline,
		StepConfigureAuthentication,
		StepConfigureDirector,
		StepDeployDirector,
		StepUploadFiles,
		StepConfigureTiles,
		StepApplyChanges,
	},
	DeleteCallback: {
		StepWaitOpsmanOnline,
		StepConfigureAuthentication,
		StepDeleteInstallation,
	},
}

// Tiler responsible for configuring Ops Manager conform a given Pattern
type Tiler struct {
	client    OpsmanClient
	logger    func(context.Context) *log.Logger
	mover     Mover
	callbacks map[callbackType][]steps.Step
}

// NewTiler creates a new Tiler needs a Mover for Tile movement
// and OpsmanClient for Ops Manager interaction
func NewTiler(client OpsmanClient, mover Mover, logger *log.Logger) *Tiler {
	l := func(ctx context.Context) *log.Logger {
		return steps.ContextLogger(ctx, logger, "[Tiler]")
	}
	return &Tiler{
		client:    client,
		mover:     mover,
		logger:    l,
		callbacks: make(map[callbackType][]steps.Step),
	}
}

// RegisterStep allows registering a Step as BuildCallback or DeleteCallback
func (t *Tiler) RegisterStep(ct callbackType, steps ...steps.Step) error {
	for _, step := range steps {
		err := validateStep(ct, step)
		if err != nil {
			return fmt.Errorf("could not register step: %v", err)
		}
		t.callbacks[ct] = append(t.callbacks[ct], step)
	}
	return nil
}

func validateStep(ct callbackType, step steps.Step) error {
	for _, name := range step.DependsOn {
		valid := false
		for _, allowedName := range allowedCallbackHooks[ct] {
			if name == allowedName {
				valid = true
			}
		}
		if !valid {
			return fmt.Errorf(
				"%s: %s may not DependOn: %s", ct, step.Name, name)
		}
	}
	return nil
}
