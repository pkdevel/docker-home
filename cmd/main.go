package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/pkdevel/docker-home/internal/pkg/docker"
	"github.com/pkdevel/docker-home/web/template/pages"
	"github.com/pkdevel/docker-home/web/template/segments"
)

func main() {
	log.Println("Setting up routes")

	// pages
	http.Handle("/", templ.Handler(pages.Index()))

	// segments
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(r.Referer())
		log.Println(url.Host)
		if err != nil {
			log.Fatal(err)
			return
		}
		hostname := strings.Split(url.Host, ":")[0]
		containers := docker.List()
		apps := ContainerApp{}
		for _, container := range containers {
			apps[container.Name] = fmt.Sprintf("http://%s:%d", hostname, container.Port)
		}
		segments.List(apps).Render(r.Context(), w)
	})

	log.Println("Starting server")
	log.Panic(http.ListenAndServe(":8080", nil))

	select {}
}

type ContainerApp map[string]string
