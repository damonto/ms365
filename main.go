package main

import (
	"github.com/damonto/msonline/cmd"
	"github.com/damonto/msonline/internal/pkg/logger"
)

func main() {
	defer logger.Sugar.Sync()

	cmd.Execute()
}
