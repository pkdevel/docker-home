package main

import (
	router "github.com/pkdevel/docker-home/internal/pkg/http"
)

func main() {
	router.SetupAndServe()
	select {}
}
