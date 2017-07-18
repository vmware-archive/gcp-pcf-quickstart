package app

import (
	"fmt"
)

type Mode string

func (m Mode) String() string {
	return string(m)
}

func (m Mode) Set(v string) error {
	m = Mode(v)
	switch m {
	case BakeImage:
	case ConfigureOpsManager:
	default:
		return fmt.Errorf("unknown mode: %s", v)
	}

	return nil
}

const (
	BakeImage           Mode = "BakeImage"
	ConfigureOpsManager      = "ConfigureOpsManager"
)
