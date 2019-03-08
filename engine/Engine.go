package engine

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"tinkodnev/utils"
)

var Router *mux.Router
var Database MemDB

func startTemplates() {
	fmt.Println("[Engine] Loading template config..")

	rawTemplateConfig := utils.MustParseJsonConfig("templates.json")
	templateConfig := TemplateConfig{
		LayoutPath:       rawTemplateConfig["layout_path"].(string),
		IncludePath:      rawTemplateConfig["include_path"].(string),
		IncludeCondition: rawTemplateConfig["include_condition"].(string),
	}

	LoadTemplates(templateConfig)
}

func startStatic() {
	fmt.Println("[Engine] Loading static files...")
	config := Router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	v, err := config.GetPathTemplate()
	if err != nil {
		fmt.Println("[Engine] Error setting up static directive: " + err.Error())
	} else {
		fmt.Println("[Engine] Started static at: " + v)
	}
}

func Load() {
	Router = new(mux.Router)
	go startTemplates()
	go startStatic()

	// Main page handler
	Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, "main.gohtml", nil)
	})

	// View page handler
	Router.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		val, found := utils.RequireU64("id", r, w);
		if !found {
			RenderTemplate(w, "error.gohtml", "Нужно указать айди")
		} else {
			RenderTemplate(w, "view.gohtml", val)
		}
	})
}

func Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("[Engine] Running at port " + port)
	err := http.ListenAndServe(":"+port, Router)
	panic(err.Error())
}
