package templatefunc

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func TrimSpaces(s string) string {
	return strings.Trim(s, "\r\n\t\v\f ")
}

func GlobDirectory(dir string, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		pattern = fmt.Sprintf("%s/*%s", dir, pattern)
		return filepath.Glob(pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			err := errors.Errorf("dir:%s filepath.Walk info is nil", dir)
			return err
		}
		if !info.IsDir() {
			path = strings.ReplaceAll(path, "\\", "/")
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func GlobFS(fsys fs.FS, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return fs.Glob(fsys, pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func getPattern(pattern string) *regexp.Regexp {
	regStr := strings.TrimLeft(pattern, ".")
	regStr = strings.ReplaceAll(regStr, ".", "\\.")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	return reg
}

func GetTemplateNames(t *template.Template) []string {
	out := make([]string, 0)
	for _, tpl := range t.Templates() {
		name := tpl.Name()
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}
