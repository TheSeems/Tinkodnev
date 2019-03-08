package engine

import (
	"fmt"
	"github.com/oxtoacart/bpool"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

/**
Specific objects
*/
var (
	templates    map[string]*template.Template
	templateBuff *bpool.BufferPool
	funcMap      template.FuncMap
)

/**
Global values
*/
var Templates template.Template
var Layouts template.Template

const baseTemplate = `
{{define "main" }}
	{{ template "base" . }}
{{ end }}
`

func SetFuncMap(fun template.FuncMap) {
	funcMap = fun
}

func LoadTemplates(config TemplateConfig) {
	if templates != nil {
		return
	}
	templates = make(map[string]*template.Template)

	Layouts = *template.Must(template.New("__layouts").
		Funcs(funcMap).
		ParseGlob(config.LayoutPath + "/" + config.IncludeCondition))
	Templates = *template.Must(template.New("__templates").
		Funcs(funcMap).
		ParseGlob(config.IncludePath + "/" + config.IncludeCondition))

	// Include, sub-templates
	includeFiles, err := filepath.Glob(config.IncludePath + "/" + config.IncludeCondition)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Layouts, main templates
	layoutFiles, err := filepath.Glob(config.LayoutPath + "/" + config.IncludeCondition)
	if err != nil {
		log.Fatal(err)
		return
	}

	mainTemplate := template.New("main").Funcs(funcMap)
	mainTemplate, err = mainTemplate.Parse(baseTemplate)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)

		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(
			templates[fileName].ParseFiles(files...))

		fmt.Println("[Templates] Loaded template: " + fileName)
	}

	templateBuff = bpool.NewBufferPool(128)
	fmt.Println("[Templates] Load finished")
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok || tmpl == nil {
		http.Error(w, fmt.Sprintf("[Templates/Error] The template %s does not exist.", name),
			http.StatusInternalServerError)
		return
	} else if templateBuff == nil {
		fmt.Println("[Templates] Buffer error for template ", tmpl)
		return
	}

	if buf := templateBuff.Get(); buf == nil {
		fmt.Println("Buf=nil")
		return
	} else {
		defer templateBuff.Put(buf)

		if buf == nil {
			fmt.Println("[Templates] Buffer error for template ", tmpl)
			return
		}

		err := tmpl.Execute(buf, data)
		if err != nil {
			fmt.Println("[Templates] Error executing template", tmpl, " => "+err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, er := buf.WriteTo(w)
		if er != nil {
			fmt.Println("[Templates] Buffer error: " + er.Error())
		}
	}
}
