package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Report is the result of scanning a module tree for exported godoc coverage.
type Report struct {
	Documented int
	Total      int
	Percent    float64
	Missing    []string
}

func measure(root string) (Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Report{}, err
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if isSourceGoFile(d.Name()) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return Report{}, err
	}

	fset := token.NewFileSet()
	var documented, total int
	var missing []string

	for _, filename := range files {
		file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return Report{}, fmt.Errorf("parse %s: %w", filename, err)
		}
		if ast.IsGenerated(file) {
			continue
		}
		rel := filename
		if r, relErr := filepath.Rel(root, filename); relErr == nil {
			rel = r
		}
		scanFile(file, rel, &documented, &total, &missing)
	}

	report := Report{
		Documented: documented,
		Total:      total,
		Missing:    missing,
	}
	if total > 0 {
		report.Percent = 100 * float64(documented) / float64(total)
	}
	return report, nil
}

func shouldSkipDir(name string) bool {
	switch name {
	case "vendor", "node_modules", "testdata", ".git", ".worktrees":
		return true
	}
	if name == "." {
		return false
	}
	return strings.HasPrefix(name, ".")
}

func isSourceGoFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

func scanFile(file *ast.File, rel string, documented, total *int, missing *[]string) {
	for _, decl := range file.Decls {
		switch x := decl.(type) {
		case *ast.FuncDecl:
			if x.Name == nil || !isExported(x.Name.Name) {
				continue
			}
			*total++
			if hasDoc(x.Doc) {
				*documented++
			} else {
				*missing = append(*missing, fmt.Sprintf("%s: func %s", filepath.ToSlash(rel), x.Name.Name))
			}
		case *ast.GenDecl:
			switch x.Tok {
			case token.TYPE, token.CONST, token.VAR:
			default:
				continue
			}
			groupDoc := hasDoc(x.Doc)
			for _, spec := range x.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if !isExported(s.Name.Name) {
						continue
					}
					*total++
					if groupDoc || hasDoc(s.Doc) {
						*documented++
					} else {
						*missing = append(*missing, fmt.Sprintf("%s: type %s", filepath.ToSlash(rel), s.Name.Name))
					}
				case *ast.ValueSpec:
					specDoc := hasDoc(s.Doc)
					for _, name := range s.Names {
						if !isExported(name.Name) {
							continue
						}
						*total++
						if groupDoc || specDoc {
							*documented++
						} else {
							*missing = append(*missing, fmt.Sprintf("%s: %s %s", filepath.ToSlash(rel), x.Tok, name.Name))
						}
					}
				}
			}
		}
	}
}

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func hasDoc(cg *ast.CommentGroup) bool {
	return cg != nil && strings.TrimSpace(cg.Text()) != ""
}
