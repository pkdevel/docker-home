package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/pkdevel/docker-home/internal/pkg/model"
	"github.com/pkdevel/docker-home/web/template/pages"
	"github.com/pkdevel/docker-home/web/template/segments"
)

func SetupAndServe() {
	go func() {
		log.Println("Setting up routes")

		// pages
		http.Handle("/", templ.Handler(pages.Index()))

		// segments
		containers := model.GetContainers()
		http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
			url, err := url.Parse(r.Referer())
			if err != nil {
				log.Panic(err)
				http.Redirect(w, r, "/500", http.StatusFound)
				return
			}
			hostname := strings.Split(url.Host, ":")[0]

			r.ParseForm()
			query := r.Form.Get("search")
			containers := containers.Find(query)
			apps := []segments.ContainerApp{}
			for _, container := range containers {
				scheme := "http"
				if container.Data.PrivatePort == 443 {
					scheme += "s"
				}

				apps = append(apps, ContainerApp{
					container.Data.Name,
					fmt.Sprintf("%s://%s:%d", scheme, hostname, container.Data.Port),
				})
			}

			segments.List(apps).Render(r.Context(), w)
		})

		// assets
		http.HandleFunc("/{file}", func(w http.ResponseWriter, r *http.Request) {
			_, err := os.Open("./assets/" + r.URL.Path[1:])
			if err != nil {
				log.Print(err)
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

		log.Println("Starting server")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()
}

type ContainerApp struct {
	name string
	url  string
}

func (c ContainerApp) Name() string {
	return c.name
}

func (c ContainerApp) URL() string {
	return c.url
}
