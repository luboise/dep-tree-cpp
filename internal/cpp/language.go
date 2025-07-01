package cpp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

type Language struct {
	Cfg *Config
}

func MakeCppLanguage(cfg *Config) (language.Language, error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(path, content)
	if err != nil {
		return nil, err
	}
	currentDir, _ := os.Getwd()
	relPath, _ := filepath.Rel(currentDir, path)
	return &language.FileInfo{
		Content: file.Statements,                    // dump the parsed statements into the FileInfo struct.
		Loc:     bytes.Count(content, []byte("\n")), // get the amount of lines of code.
		Size:    len(content),                       // get the size of the file in bytes.
		AbsPath: path,                               // provide its absolute path.
		RelPath: relPath,                            // provide the path relative to the current dir.
	}, nil
}

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	var result language.ImportsResult

	for _, statement := range file.Content.([]Statement) {
		fmt.Println("Found include: ", statement)

		if statement.Quoted != nil {
			fmt.Println("Found quoted: ", statement.Quoted)
			var path = *&statement.Quoted.IncludedFile

			if strings.HasPrefix(path, ".") {
				fmt.Println("FOUND A RELATIVE IMPORT")
				result.Imports = append(result.Imports, language.ImportEntry{
					Symbols: []string{path},
					AbsPath: filepath.Join(filepath.Dir(file.AbsPath), path),
				})
			} else {
				fmt.Println("FOUND AN ABSOLUTE IMPORT")
				if l.Cfg == nil {
					fmt.Println("NO CONFIG FOUND. SKIPPING ABSOLUTE IMPORT.")
					continue
				}

				var found = false
				for _, includePath := range l.Cfg.IncludePaths {
					var absPath = filepath.Join(includePath, path)

					if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
						continue
					}

					fmt.Println("Found included path ", path, "at absolute location ", absPath)
					found = true
					break
				}
				if !found {
					fmt.Println("Unable to find included path: ", path)
					continue
				}

				result.Imports = append(result.Imports, language.ImportEntry{
					// TODO: Get the symbols from the other file instead of using the header file
					Symbols: []string{path},
					AbsPath: filepath.Join(filepath.Dir(file.AbsPath), path),
				})
			}
		} else if statement.Angled != nil {
			var path = statement.Angled.IncludedFile

			fmt.Println("Found angled")
			result.Imports = append(result.Imports, language.ImportEntry{
				// TODO: Get the symbols from the other file instead of using the header file
				Symbols: []string{path},
				AbsPath: filepath.Join(filepath.Dir(file.AbsPath), path),
			})
		}
	}

	return &result, nil
}

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	var result language.ExportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Quoted != nil {
			var header = statement.Quoted.IncludedFile

			result.Exports = append(result.Exports, language.ExportEntry{
				Symbols: []language.ExportSymbol{{Original: header}},
				AbsPath: file.AbsPath,
			})
		} else if statement.Angled != nil {
			var header = statement.Angled.IncludedFile

			result.Exports = append(result.Exports, language.ExportEntry{
				Symbols: []language.ExportSymbol{{Original: header}},
				AbsPath: file.AbsPath,
			})

		}
	}

	return &result, nil
}

var Extensions = []string{"h", "cpp", "cppm", "ixx"}
