package router

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/pkdevel/docker-home/internal/pkg/model"
	"github.com/pkdevel/docker-home/web/template/pages"
	"github.com/pkdevel/docker-home/web/template/segments"
)

func SetupAndServe() {
	slog.Info("Setting up routes")

	// pages
	http.Handle("/{$}", templ.Handler(pages.Index()))

	// segments
	endpoints := model.GetEndpoints()
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("dhcq-search")
		apps := []*segments.ContainerApp{}
		for _, endpoint := range endpoints.Find(query) {
			if len(endpoint.Links) == 0 {
				continue
			}
			for _, link := range endpoint.Links {
				apps = append(apps, &segments.ContainerApp{
					Name: endpoint.ID,
					URL:  link,
				})
			}
		}
		segments.Containers(apps).Render(r.Context(), w)
	})

	// assets
	http.HandleFunc("/{file}", func(w http.ResponseWriter, r *http.Request) {
		filename := fmt.Sprintf("./assets/%s", r.PathValue("file"))
		_, err := os.Open(filename)
		if err != nil {
			slog.Error(err.Error())
			http.Redirect(w, r, "/404", http.StatusFound)
			return
		}
		http.ServeFile(w, r, filename)
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

	slog.Info("Starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
