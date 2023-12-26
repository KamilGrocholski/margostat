package ssg

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"scraper/config"
	"scraper/database"
)

type WorldPageData struct {
	Worlds        []string
	SelectedWorld string
	CountResults  []database.CountResult
	GeneratedAt   time.Time
}

const (
	TEMPLATE_DIR_PATH string = "templates"
)

func ComposePathToStaticFile(name string, ext string) string {
	return fmt.Sprintf("%s/%s.%s", config.STATIC_DIR_PATH, name, ext)
}

func ComposePathToTemplateFile(name string, ext string) string {
	return fmt.Sprintf("%s/%s.%s", TEMPLATE_DIR_PATH, name, ext)
}

func GenerateAndWriteHtmlPageFileToStatic(filename string, data *WorldPageData) error {
	destPath := ComposePathToStaticFile(filename, "html")

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %v", destPath, err)
	}
	defer file.Close()

	templatePath := ComposePathToTemplateFile("layout", "html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing HTML template: %v", err)
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing HTML template: %v", err)
	}

	return nil
}
