package main

import (
	"github.com/damonto/msonline-webapi/cmd"
	"github.com/damonto/msonline-webapi/internal/pkg/logger"
)

func main() {
	defer logger.Sugar.Sync()

	cmd.Execute()
}
