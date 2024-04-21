package server

import (
	"log"
	"net/http"
	"text/template"
)

type PageData struct {
	Title string
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	templateFile := "./views/index.html"
	tmpl := template.Must(template.ParseFiles(templateFile))

	pageData := PageData{
		Title: "Recipe Refiner",
	}

	err := tmpl.Execute(w, pageData)

	if err != nil {
		log.Fatalf("Error applying template: %v", err)
	}
}

func RunServer(host string, port string) error {
	address := host + ":" + port

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", getRoot)

	return http.ListenAndServe(address, mux)
}
