package server

import (
	"github.com/simplepki/server/config"
	"github.com/simplepki/core"
)

func InitializeCA() error {
	if config.IsCAEnabled() && config.ShouldOverwriteCA() {
		//new ca

	} else {
		// error if ca is not whats expected
	}

	return nil
}
