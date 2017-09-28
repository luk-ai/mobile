package bind

import (
	"go/build"
	"io"
	"path/filepath"
	"text/template"
)

type NativeMeta struct {
	Libs []string
}

type Writer func(string) (io.Writer, func(), error)

// GenerateJavaSupport generates the java support files.
func GenerateJavaSupport(mobilePkg *build.Package, meta NativeMeta, dir string, writer Writer) error {
	repo := filepath.Clean(filepath.Join(mobilePkg.Dir, "..")) // golang.org/x/mobile directory.
	files := []string{"Seq.java", "LoadJNI.java"}
	var fullFiles []string
	for _, file := range files {
		fullFiles = append(fullFiles, filepath.Join(repo, "bind/java", file))
	}
	templates, err := template.ParseFiles(fullFiles...)
	if err != nil {
		return err
	}
	for _, javaFile := range files {
		w, closer, err := writer(filepath.Join(dir, javaFile))
		if err != nil {
			return err
		}
		defer closer()
		if err := templates.ExecuteTemplate(w, javaFile, meta); err != nil {
			return err
		}
	}
	return nil
}
