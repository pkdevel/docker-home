package main

import (
	"log/slog"

	router "github.com/pkdevel/docker-home/internal/pkg/http"
	"github.com/pkdevel/docker-home/internal/pkg/persistence"
	"github.com/pkdevel/docker-home/internal/pkg/task"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	persistence.Init()
	defer persistence.Close()
	task.StartImporter()
	go router.SetupAndServe()

	select {}
}
