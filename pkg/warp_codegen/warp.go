package warp

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type ITemplate interface {
	Generate() error
}

type Template struct {
	Elems []string
	//Ifaces  []any
	FuncMap template.FuncMap
}

const serverTemplateFile = "./lib/server_grpc/src/server.rs"

func NewTemplate(cleanBuild bool) ITemplate {
	elems := []string{
		"server_rs",
	}
	var isCleanBuild = func() bool {
		return cleanBuild
	}
	fMap := template.FuncMap{
		"IsCleanBuild": isCleanBuild,
	}
	return &Template{
		Elems:   elems,
		FuncMap: fMap,
	}
}
func (t *Template) Generate() error {
	file, _ := os.Create(serverTemplateFile)

	// nolint:errcheck
	defer file.Close()

	pattern, _ := filepath.Abs("pkg/warp_codegen/server.gotml") // .gotmpl is used because of IDE's supports only :D

	tmpl := template.Must(template.New("").Funcs(t.FuncMap).ParseFiles(pattern))
	//var wg sync.WaitGroup
	//for i := range t.Ifaces {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	err := tmpl.ExecuteTemplate(file, t.Elems[0], nil) // first arg is output, second is the data we want to pass to this config. It could also be nil.
	if err != nil {
		return fmt.Errorf("template generation failed: %v", err)
	}
	//if err != nil {
	//	logger.Log().Errorf("An error occurred %v", err)
	//	return
	//}
	//}()
	//wg.Wait()
	//}

	return nil
}

func RemoveTargerFolder() error {
	if err := os.RemoveAll("./lib/server_grpc/target"); err != nil {
		return fmt.Errorf("error while removing target folder: %v", err)
	}
	return nil
}
