package main

import (
	"github.com/damonto/ms365/cmd"
	"github.com/damonto/ms365/internal/pkg/logger"
)

func main() {
	defer logger.Sugar.Sync()

	cmd.Execute()
}
