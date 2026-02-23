package main

import (
	"fmt"

	"github.com/fluxergo/fluxergo/fluxer"
)

// leetify stats auto messaging and commands handler

func getStatsHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	profile, err := getLeetifyStats(args[1])

	if err != nil {
		return fmt.Errorf("error getting stats: %w", err)
	}

	stat := profile
	fmt.Println("requested stat:", stat)

	return nil
}
