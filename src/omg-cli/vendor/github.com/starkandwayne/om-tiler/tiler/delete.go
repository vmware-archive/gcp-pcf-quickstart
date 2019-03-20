package tiler

func (c *Tiler) Delete() error {
	err := c.client.ConfigureAuthentication()
	if err != nil {
		return err
	}

	return c.client.DeleteInstallation()
}
