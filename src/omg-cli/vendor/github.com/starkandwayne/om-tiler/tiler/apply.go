package tiler

func (c *Tiler) Apply() error {
	return c.client.ApplyChanges()
}
