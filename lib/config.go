package lib

/*var (*/
//datafileFlag = flag.String("data", "", "Datafile")
//tmplFlag     = flag.String("tmpl", "", "Template")

//tmpl []byte
//data map[interface{}]interface{}
//)

//func init() {
//flag.Parse()

//stat, _ := os.Stdin.Stat()
//if (stat.Mode()&os.ModeNamedPipe == 0) && *tmplFlag == "" {
//panic("No template ")
//}

//var (
//source *os.File
//)

//if *tmplFlag == "-" {
//source = os.Stdin
//defer source.Close()
//} else {
//files := strings.Split(*tmplFlag, ",")
//sources := make([]*os.File, len(files))
//for idx, file := range files {
//source, err := os.Open(file)
//if err != nil {
//panic(err)
//}
//sources[idx] = source
//defer source.Close()
//}
//for _, source := range sources {
//aTmpl, err := ioutil.ReadAll(source)
//if err != nil {
//panic(err)
//}
//tmpl = append(tmpl, aTmpl...)
//}
//}

//if *datafileFlag != "" {
//files := strings.Split(*datafileFlag, ",")
//out, err := parseAll(files)
//if err != nil {
//panic(err)
//}
//data = out
//}
//}

//func parseAll(filepaths []string) (map[interface{}]interface{}, error) {
//maps := make([]map[interface{}]interface{}, len(filepaths))

//for idx, filepath := range filepaths {
//d, err := parse(filepath)
//if err != nil {
//return nil, err
//}
//maps[idx] = d
//}

//output := make(map[interface{}]interface{})
//for _, d := range maps {
//output = merge(output, d)
//}

//return output, nil
//}

//func parse(filepath string) (map[interface{}]interface{}, error) {
//var d map[interface{}]interface{}

//dataBytes, err := ioutil.ReadFile(filepath)
//if err != nil {
//return nil, err
//}

//err = yaml.Unmarshal(dataBytes, &d)
//if err != nil {
//return nil, err
//}

//return d, nil
//}

//// merge takes two maps and merges them. on collision b overwrites a
//func merge(a, b map[interface{}]interface{}) map[interface{}]interface{} {
//output := make(map[interface{}]interface{})

//for k, v := range a {
//output[k] = v
//}

//for k, v := range b {
//output[k] = v
//}

//return output
//}

//func main() {
//t := template.Must(template.New(*tmplFlag).Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))
//if err := t.Execute(os.Stdout, data); err != nil {
//panic(err)
//}
/*}*/
