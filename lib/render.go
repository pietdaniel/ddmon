package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/hcl2/hclwrite"
	yaml "gopkg.in/yaml.v2"
)

func render(tplPaths, dataPaths []string, filename, outputPath string) {
	data, err := parseAll(dataPaths)
	if err != nil {
		panic(err)
	}

	tmpl := loadTemplates(tplPaths)

	t := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))

	// todo needs to write to output fath and write to name

	outputFile := fmt.Sprintf("%s/%s.tf", outputPath, filename)
	log.Printf("Rendering data into template for tf file: %s", outputFile)

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, data); err != nil {
		log.Printf("Failed to render template: %v %v %s", tplPaths, dataPaths, err)
	}
	// todo extract to function

	f.Close()

	log.Printf("Formatting file %s", outputFile)

	r, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	_, syntaxDiags := hclsyntax.ParseConfig(src, outputFile, hcl.Pos{Line: 1, Column: 1})

	if syntaxDiags.HasErrors() {
		log.Printf("Failed to parse tf file %s: %v", outputFile, syntaxDiags)
	}

	formatted := hclwrite.Format(src)

	if !bytes.Equal(src, formatted) {
		_, err := bytesDiff(src, formatted, outputFile)
		if err != nil {
			log.Fatal(err)
		}
		// We could print the diff returned from bytesDiff here if need be
	}

	r.Close()

	os.Remove(outputFile)

	f, err = os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(formatted)
	if err != nil {
		log.Fatal(err)
	}
}

func bytesDiff(b1, b2 []byte, path string) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "--label=old/"+path, "--label=new/"+path, "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}

func loadTemplates(tplPaths []string) []byte {
	var tmpl []byte
	for _, tplPath := range tplPaths {
		log.Printf("Loading template: %s", tplPath)
		tpl, err := os.Open(tplPath)
		defer tpl.Close()
		if err != nil {
			panic(err)
		}

		aTmpl, err := ioutil.ReadAll(tpl)
		if err != nil {
			panic(err)
		}
		tmpl = append(tmpl, aTmpl...)
	}
	return tmpl
}

// merge takes two maps and merges them. on collision b overwrites a
func merge(a, b map[interface{}]interface{}) map[interface{}]interface{} {
	output := make(map[interface{}]interface{})

	for k, v := range a {
		output[k] = v
	}

	for k, v := range b {
		output[k] = v
	}

	return output
}

// parseAll tkaes a list of filepaths and returns a map of keys to values
// collisions are overwritten by later entries in the list
func parseAll(filepaths []string) (map[interface{}]interface{}, error) {
	maps := make([]map[interface{}]interface{}, len(filepaths))

	for idx, filepath := range filepaths {
		log.Printf("Parsing data from %v", filepath)
		d, err := parse(filepath)
		if err != nil {
			return nil, err
		}
		maps[idx] = d
	}

	output := make(map[interface{}]interface{})
	for _, d := range maps {
		output = merge(output, d)
	}

	return output, nil
}

// parse takes a filepath for a yaml file and returns a map of keys and values
func parse(filepath string) (map[interface{}]interface{}, error) {
	var d map[interface{}]interface{}

	dataBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dataBytes, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
