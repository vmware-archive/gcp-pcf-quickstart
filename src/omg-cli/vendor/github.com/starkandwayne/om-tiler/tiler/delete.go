package tiler

import "context"

func (c *Tiler) Delete(ctx context.Context) error {
	err := c.client.ConfigureAuthentication(ctx)
	if err != nil {
		return err
	}

	return c.client.DeleteInstallation(ctx)
}
