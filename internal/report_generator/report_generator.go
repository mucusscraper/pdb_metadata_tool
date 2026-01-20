package reportgenerator

import (
	"fmt"
	"html/template"
	"os"

	getdata "github.com/mucusscraper/pdb_metadata_tool/internal/get_data"
)

type Report struct {
	PreReport []getdata.PreReport
	Grouped   bool
}

func GenerateHTML(filename string, data Report) error {
	templ, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return err
	}
	filepath := fmt.Sprintf("reports/%v.html", filename)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return templ.Execute(file, data)
}
