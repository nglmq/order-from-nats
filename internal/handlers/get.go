package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func TemplateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/main.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Println(r.Method, r.URL.Path)
	}
}
