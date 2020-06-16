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
	templateDir := cmd.Flags().Lookup("template-dir").Value.String()
	log.Printf("Generating from: %s to: %s. Template Dir: %s", sourceDir, targetDir, templateDir)

	mustGetPath(targetDir)
	mustGetPath(sourceDir)
	mustGetPath(templateDir)

	// Find data directory
	dataDir := mustGetPath(fmt.Sprintf("%s/data", sourceDir))

	// Look for common.yaml in base data directory
	dataCommonFilePath := mustGetPath(fmt.Sprintf("%s/common.yaml", dataDir))

	// Aggregate all namepsaces "data/$NAMESPACE/$GROUP/monitor-a.yaml
	namespaces := getFolderNames(dataDir)
	for _, namespace := range namespaces {
		log.Printf("Generating on namespace: %v", namespace)
		namespaceDir := mustGetPath(fmt.Sprintf("%s/%s", dataDir, namespace))

		// Look for namespace common.yaml
		namespaceCommonFilePath := mustGetPath(fmt.Sprintf("%s/common.yaml", namespaceDir))

		groups := getFolderNames(namespaceDir)
		for _, group := range groups {
			log.Printf("Generating on group: %s", group)

			// Look for group common.yaml
			groupDir := mustGetPath(fmt.Sprintf("%s/%s", namespaceDir, group))
			groupCommonFilePath := mustGetPath(fmt.Sprintf("%v/common.yaml", groupDir))

			commonDataFiles := []string{dataCommonFilePath, namespaceCommonFilePath, groupCommonFilePath}

			monitorDataFiles := getMonitorDataFilesPaths(groupDir)
			for _, f := range monitorDataFiles {
				filePaths := append(commonDataFiles, f)
				monitor, err := getMonitorDataFile(f)
				if err != nil {
					panic(err)
				}

				templates := getListOfTemplates(templateDir, monitor, group, namespace)

				filename := fmt.Sprintf("%s-%s", group, monitor)

				render(templates, filePaths, filename, targetDir)
			}
		}
	}
}

// mustGetFilePath gets a path to a file, if the file does not exist it fatally exits
func mustGetPath(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		log.Printf("Could not find path: %s", path)
		log.Fatal(err)
	}
	return path
}

// getListOfTemplates will return the list of templates file paths as strings for the given monitor
// the heuristic is as follows:
//   if there is a template with the $MONITOR_NAME.tpl return [$MONITOR_NAME.tpl, base.tpl]
//   if there is a template with the $GROUP.tpl return [$GROUP.tpl, base.tpl]
//   if there is a template with the $NAMESPACE.tpl return [$NAMESPACE.tpl, base.tpl]
//   else return [default.tpl, base.tpl]
func getListOfTemplates(templateDir, monitor, group, namespace string) []string {
	baseTplFilePath := mustGetPath(fmt.Sprintf("%s/%s", templateDir, "base.tpl"))
	defaultTplFilePath := mustGetPath(fmt.Sprintf("%s/%s", templateDir, "default.tpl"))

	templates := []string{baseTplFilePath}

	// Template mapping
	if tpl, ok := checkForTplFile(monitor, templateDir); ok {
		templates = append(templates, tpl)
	} else if tpl, ok := checkForTplFile(group, templateDir); ok {
		templates = append(templates, tpl)
	} else if tpl, ok := checkForTplFile(namespace, templateDir); ok {
		templates = append(templates, tpl)
	} else {
		templates = []string{baseTplFilePath, defaultTplFilePath}
	}

	return templates
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
