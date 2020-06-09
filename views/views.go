package views

import (
	"net/http"
	"path/filepath"
	"text/template"

	"profile.com/context"
)

const (
	fileExt    = ".html"
	filePrefix = "views/"
)

// Views defines the shape of the views
type Views struct {
	t      *template.Template
	layout string
}

func layoutFiles() []string {
	files, err := filepath.Glob("views/layout/*")
	if err != nil {
		panic(err)
	}
	return files
}

func handleExts(files []string) {
	for i, f := range files {
		files[i] = filePrefix + f + fileExt
	}
}

// NewView returns the views struct
func NewView(layout, file string) *Views {
	files := []string{file}
	handleExts(files)
	lf := layoutFiles()
	files = append(files, lf...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &Views{
		t:      t,
		layout: layout,
	}
}

// Render renders the page
func (v *Views) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var vd Data
	if d, t := data.(Data); t {
		vd = d
	} else {
		vd = Data{
			Yield: data,
		}
	}

	user := context.GetUserFromContext(r.Context())
	vd.User = user
	if err := v.t.ExecuteTemplate(w, v.layout, vd); err != nil {
		panic(err)
	}
}
