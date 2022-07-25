package main

import "github.com/fishykins/aero/pkg/logging"

var Log *logging.Logger

func main() {
	Log = logging.New()
	Log.Info("Hello, world!")
	Log.InfoWith("Hello, world!", map[string]interface{}{"name": "fishykins"})
}
