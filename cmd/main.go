package main

import (
	router "github.com/pkdevel/docker-home/internal/pkg/http"
	"github.com/pkdevel/docker-home/internal/pkg/task"
)

func main() {
	router.SetupAndServe()
	task.StartImporter()
	select {}
}
