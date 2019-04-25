package tiler

import (
	"context"

	"github.com/starkandwayne/om-tiler/steps"
)

// Delete will perform a series of steps to delete all installed Tiles.
func (t *Tiler) Delete(ctx context.Context) error {
	s := []steps.Step{
		t.stepPollTillOnline(),
		t.stepConfigureAuthentication(),
		t.stepDeleteInstallation(),
	}
	s = append(s, t.callbacks[DeleteCallback]...)

	return steps.Run(ctx, s, t.logger(ctx))
}
