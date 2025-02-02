package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/stevelowery/driftless/internal/cmd"
)

func main() {

	ctx := context.Background()

	if err := cmd.New().ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
