package cpp

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
)

type Language struct{}

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
		if statement.Include != nil {
			if statement.Include.Quoted != nil {
				result.Imports = append(result.Imports, language.ImportEntry{
					// TODO: Get the symbols from the other file instead of using the header file
					Symbols: []string{*statement.Include.Quoted},
					AbsPath: filepath.Join(filepath.Dir(file.AbsPath), *statement.Include.Quoted),
				})
			}
		} else if statement.Include.Angled != nil {
			result.Imports = append(result.Imports, language.ImportEntry{
				// TODO: Get the symbols from the other file instead of using the header file
				Symbols: []string{*statement.Include.Angled},
				AbsPath: filepath.Join(filepath.Dir(file.AbsPath), *statement.Include.Angled),
			})
		}
	}

	return &result, nil
}

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	var result language.ExportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Include != nil {
			if statement.Include.Quoted != nil {

				result.Exports = append(result.Exports, language.ExportEntry{
					Symbols: []language.ExportSymbol{{Original: *statement.Include.Angled}},
					AbsPath: file.AbsPath,
				})

			}
		}
	}

	return &result, nil
}

var Extensions = []string{"h", "cpp", "cppm", "ixx"}
