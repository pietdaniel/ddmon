package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

// Generate genearates TF files
func Generate(cmd *cobra.Command, args []string) {
	// Get arguments
	targetDir := cmd.Flags().Lookup("target-dir").Value.String()
	sourceDir := cmd.Flags().Lookup("source-dir").Value.String()
	log.Printf("Generating from: %s to: %s", sourceDir, targetDir)

	deleteFilesFromOutputDirectory(targetDir)

	// Find data directory
	dataDir := fmt.Sprintf("%s/data", sourceDir)
	if _, err := os.Stat(dataDir); os.IsExist(err) {
		log.Fatalf("Could not find data directory at: %s", dataDir)
	}

	// Find template directory
	templateDir := fmt.Sprintf("%s/templates", sourceDir)
	if _, err := os.Stat(templateDir); os.IsExist(err) {
		log.Fatalf("Could not find data directory at: %s", templateDir)
	}

	// Look for common.yaml in base data directory
	dataCommonFilePath := fmt.Sprintf("%s/common.yaml", dataDir)
	_, err := os.Stat(dataCommonFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Aggregate all namepsaces "data/$NAMESPACE/$GROUP/monitor-a.yaml
	namespaces := getFolderNames(dataDir)
	for _, namespace := range namespaces {
		log.Printf("Generating on namespace: %v", namespace)
		namespaceDir := fmt.Sprintf("%s/%s", dataDir, namespace)

		// Look for namespace common.yaml
		namespaceCommonFilePath := fmt.Sprintf("%s/common.yaml", namespaceDir)
		_, err := os.Stat(namespaceCommonFilePath)
		if err != nil {
			log.Fatal(err)
		}

		groups := getFolderNames(namespaceDir)
		for _, group := range groups {
			log.Printf("Generating on group: %s", group)

			groupDir := fmt.Sprintf("%s/%s", namespaceDir, group)
			// Look for namespace common.yaml
			groupCommonFilePath := fmt.Sprintf("%v/common.yaml", groupDir)
			_, err := os.Stat(groupCommonFilePath)
			if err != nil {
				log.Fatal(err)
			}

			templatesDir := fmt.Sprintf("%s/templates", sourceDir)
			baseTplFilePath := fmt.Sprintf("%s/%s", templatesDir, "base.tpl")
			_, err = os.Stat(baseTplFilePath)
			if err != nil {
				log.Fatal(err)
			}

			defaultTplFilePath := fmt.Sprintf("%s/%s", templatesDir, "default.tpl")
			_, err = os.Stat(defaultTplFilePath)
			if err != nil {
				log.Fatal(err)
			}

			commonDataFiles := []string{dataCommonFilePath, namespaceCommonFilePath, groupCommonFilePath}

			monitorDataFiles := getMonitorDataFilesPaths(groupDir)
			for _, f := range monitorDataFiles {
				filePaths := append(commonDataFiles, f)
				monitor, err := getMonitorDataFile(f)
				if err != nil {
					panic(err)
				}

				templates := []string{baseTplFilePath}

				if tpl, ok := checkForTplFile(monitor, sourceDir); ok {
					templates = append(templates, tpl)
				} else if tpl, ok := checkForTplFile(group, sourceDir); ok {
					templates = append(templates, tpl)
				} else if tpl, ok := checkForTplFile(namespace, sourceDir); ok {
					templates = append(templates, tpl)
				} else {
					templates = []string{baseTplFilePath, defaultTplFilePath}
				}

				render(templates, filePaths, monitor, targetDir)
			}
		}
	}
}

func checkForTplFile(name string, sourceDir string) (string, bool) {
	tplFilePath := fmt.Sprintf("%s/templates/%s.tpl", sourceDir, name)
	_, err := os.Stat(tplFilePath)
	if err != nil {
		log.Printf("Could not find template %s", tplFilePath)
		return tplFilePath, false
	}
	return tplFilePath, true
}

func getFolderNames(dir string) []string {
	folderNames := []string{}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			folderNames = append(folderNames, fileInfo.Name())
		}
	}
	return folderNames
}

func getMonitorDataFile(monitorPath string) (string, error) {
	log.Print("Attempting to parse monitor data file name")
	re := regexp.MustCompile(`/(?P<MonitorName>[a-zA-Z0-9\-]*).yaml$`)
	matched := re.FindStringSubmatch(monitorPath)
	if len(matched) < 1 {
		return "", fmt.Errorf("Failed to parse monitor name from path %s", monitorPath)
	}
	return matched[1], nil
}

func getMonitorDataFilesPaths(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != "common.yaml" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return files
}

func deleteFilesFromOutputDirectory(outputPath string) {
	log.Printf("Clearing files from %s", outputPath)
	err := filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			log.Printf("Remove file: %s", path)
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
