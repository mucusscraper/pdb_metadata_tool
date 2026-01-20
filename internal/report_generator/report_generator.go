package reportgenerator

import (
	"embed"
	"fmt"
	"html/template"
	"os"

	getdata "github.com/mucusscraper/pdb_metadata_tool/internal/get_data"
)

type Report struct {
	PreReport []getdata.PreReport
	Grouped   bool
}

//go:embed templates/*.html
var templateFS embed.FS

func GenerateHTML(filename string, data Report) error {
	templ, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		return err
	}
	filepath := fmt.Sprintf("reports/%v.html", filename)
	err = os.MkdirAll("reports", 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return templ.Execute(file, data)
}
