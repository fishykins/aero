package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fishykins/aero/internal"
	"go.uber.org/zap"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "bin/logs/engine.txt"}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()

	ticker := time.NewTicker(500 * time.Millisecond)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	app := internal.App{
		Logger: sugar,
	}
	app.Init()

	reader := bufio.NewReader(os.Stdin)
	go func() {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		app.ConsoleInput(text)
	}()

	for {
		select {
		case t := <-ticker.C:
			app.Run(&t)
		case <-shutdown:
			app.Shutdown()
			os.Exit(1)
		}
	}
}
