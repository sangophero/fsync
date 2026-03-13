package main

import (
	"fsync/cmd"

	"github.com/rs/zerolog/log"
)

func main() {
	rootCommandHandler := cmd.NewRootCommandHandler()

	if err := rootCommandHandler.Command.Execute(); err != nil {
		log.Err(err).Msg("execution error")
		return
	}
}
