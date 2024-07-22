package router

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

func SetupAndServe() {
	log.Println("Setting up routes")

	// pages
	http.Handle("/", templ.Handler(pages.Index()))

	// segments
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(r.Referer())
		if err != nil {
			log.Fatal(err)
			http.Redirect(w, r, "/500", http.StatusFound)
			return
		}
		hostname := strings.Split(url.Host, ":")[0]
		containers := docker.List()
		apps := ContainerApp{}
		for _, container := range containers {
			scheme := "http"
			if container.PrivatePort == 443 {
				scheme += "s"
			}
			apps[container.Name] = fmt.Sprintf("%s://%s:%d", scheme, hostname, container.Port)
		}
		segments.List(apps).Render(r.Context(), w)
	})

	// assets
	http.HandleFunc("/{file}", func(w http.ResponseWriter, r *http.Request) {
		_, err := os.Open("./assets/" + r.URL.Path[1:])
		if err != nil {
			log.Fatal(err)
			http.Redirect(w, r, "/404", http.StatusFound)
			return
		}
		http.ServeFile(w, r, "./assets/"+r.URL.Path[1:])
	})

	// errors
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		pages.NotFound().Render(r.Context(), w)
	})
	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error().Render(r.Context(), w)
	})

	go func() {
		log.Println("Starting server")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()
}

type ContainerApp map[string]string
