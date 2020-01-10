package main

import (
	"github.com/damonto/office365/cmd"
	"github.com/damonto/office365/internal/pkg/logger"
)

func main() {
	defer logger.Sugar.Sync()

	cmd.Execute()
}
