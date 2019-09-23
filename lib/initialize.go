package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gobuffalo/packr"
)

// Run x
func Run(args []string) {
	log.Println("Called intialize")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)

	log.Println("Getting templates")
	box := packr.NewBox("./templates")

	// make data directory
	log.Println("Creating data directory")
	dataPath := fmt.Sprintf("%v/%v", dir, "data")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.Mkdir(dataPath, os.ModePerm)
	}

	// make templates directory
	log.Println("Creating templates directory")
	templatesPath := fmt.Sprintf("%v/%v", dir, "templates")
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		os.Mkdir(templatesPath, os.ModePerm)
	}

	log.Println("Creating templates files")
	// make resources/README
	readme, err := box.FindString("tplREADME.md")
	if err != nil {
		log.Fatal(err)
	}

	templateReadmePath := fmt.Sprintf("%v/%v", dir, "templates/README.md")
	err = ioutil.WriteFile(templateReadmePath, []byte(readme), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// make resources/templates/base.tpl
	base, err := box.FindString("base.tpl")
	if err != nil {
		log.Fatal(err)
	}

	baseTplPath := fmt.Sprintf("%v/%v", dir, "templates/base.tpl")
	err = ioutil.WriteFile(baseTplPath, []byte(base), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// make resources/templates/default.tpl
	defaultTpl, err := box.FindString("default.tpl")
	if err != nil {
		log.Fatal(err)
	}

	defaultTplPath := fmt.Sprintf("%v/%v", dir, "templates/default.tpl")
	err = ioutil.WriteFile(defaultTplPath, []byte(defaultTpl), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// make monitors directory for terraform output
	log.Println("Creating monitors directory")
	monitorsPath := fmt.Sprintf("%v/%v", dir, "monitors")
	if _, err := os.Stat(monitorsPath); os.IsNotExist(err) {
		os.Mkdir(monitorsPath, os.ModePerm)
	}

	// make root README.md
	rootReadme, err := box.FindString("rootREADME.md")
	if err != nil {
		log.Fatal(err)
	}

	rootReadmePath := fmt.Sprintf("%v/%v", dir, "README.md")
	err = ioutil.WriteFile(rootReadmePath, []byte(rootReadme), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// make output folder
	log.Println("Creating output directory")
	outputPath := fmt.Sprintf("%v/%v", dir, "output")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}
}
