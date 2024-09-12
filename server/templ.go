package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Templates struct {
	Tmpl map[string]string
	mu   sync.Mutex
}

func CheckTmplPath(path string) (bool, error) {
	// Convert the path to an absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	// Check if the file exists and is not a directory
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("path does not exist: %s", absPath)
	}
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("path is a directory, not a file: %s", absPath)
	}
	fmt.Println("path:", absPath)
	return true, nil
}

func (t *Templates) AddTemplate(name, path string) error {
	// Use the checkPath function to validate the path
	valid, err := CheckTmplPath(path)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("invalid template path: %s", path)
	}

	t.Tmpl[name] = path
	return nil
}

func HandleTemplate(file string) http.HandlerFunc {
	var (
		tmpl *template.Template
		err  error
		init sync.Once
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tmpl, err = template.ParseFiles(file)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
