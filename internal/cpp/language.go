package cpp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

type Language struct {
	Cfg                 *Config
	AllowedSTLFilepaths []string
}

func MakeCppLanguage(cfg *Config) (language.Language, error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) GetIncludePath(path string) (includePath string, recursive bool, err error) {
	if l.Cfg != nil {
		for _, includePath := range l.Cfg.RecursiveIncludePaths {
			if strings.HasPrefix(path, includePath) {
				return includePath, true, nil
			}
		}

		for _, includePath := range l.Cfg.NonRecursiveIncludePaths {
			// If file is in stl
			if strings.HasPrefix(path, includePath) {
				// If file hasn't been included from a non-stl filepath, then skip the file
				if !slices.Contains(l.AllowedSTLFilepaths, path) {
					return includePath, false, nil
				}
				break
			}
		}
	}

	return "", false, os.ErrNotExist
}

func (l *Language) GetABSPath(includedPath string) (found bool, absPath string, isRecursive bool) {
	if l.Cfg == nil {
		return false, "", false
	}

	for _, includePath := range l.Cfg.RecursiveIncludePaths {
		var innerAbsPath = filepath.Join(includePath, includedPath)

		if _, err := os.Stat(innerAbsPath); errors.Is(err, os.ErrNotExist) {
			continue
		}

		fmt.Println("Found included path ", includePath, "at absolute location ", innerAbsPath)
		absPath = filepath.Clean(innerAbsPath)
		isRecursive = true
		found = true

		break
	}

	if !found {
		for _, includePath := range l.Cfg.NonRecursiveIncludePaths {
			var innerAbsPath = filepath.Join(includePath, includedPath)

			if _, err := os.Stat(innerAbsPath); errors.Is(err, os.ErrNotExist) {
				continue
			}

			fmt.Println("Found included path ", includePath, "at absolute location ", innerAbsPath)
			absPath = filepath.Clean(innerAbsPath)
			isRecursive = false
			found = true
			break
		}
	}

	return found, filepath.Clean(absPath), isRecursive

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
	fmt.Println("found a banchod ", file)
	var result language.ImportsResult

	_, isRecursive, err := l.GetIncludePath(file.AbsPath)
	if err != nil {
		fmt.Println(`Error occurred (Skipping `, file.AbsPath, `): `, err)
		return &result, nil
		// If the file is from the STL and isn't on the exception list, skip it
	} else if !isRecursive && !slices.Contains(l.AllowedSTLFilepaths, file.AbsPath) {
		return &result, nil
	}

	for _, statement := range file.Content.([]Statement) {
		fmt.Println("Found include: ", statement)

		var includePath string = ""

		if statement.Quoted != nil {
			includePath = statement.Quoted.IncludedFile
			fmt.Println("Found quoted: ", includePath)
		} else if statement.Angled != nil {
			includePath = statement.Angled.IncludedFile
			fmt.Println("Found angled: ", includePath)
		} else {
			continue
		}

		// If it's a relative include, add the filepaths together and add it as an Import
		if strings.HasPrefix(includePath, ".") {
			var absInclude = filepath.Clean(filepath.Join(filepath.Dir(file.AbsPath), includePath))
			result.Imports = append(result.Imports, language.ImportEntry{
				Symbols: []string{absInclude},
				AbsPath: absInclude,
			})
			continue
		}

		// Otherwise, get the abs path of the include
		found, absPath, isRecursive := l.GetABSPath(includePath)

		if !found {
			fmt.Println("Unable to find included path: ", includePath)
			continue
		}

		if isRecursive {
			l.AllowedSTLFilepaths = append(l.AllowedSTLFilepaths, absPath)
		}

		result.Imports = append(result.Imports, language.ImportEntry{
			Symbols: []string{absPath},
			AbsPath: absPath,
		})

		/*
			 else if statement.Angled != nil {
					var path = statement.Angled.IncludedFile

					var absPath string = ""
					for _, includePath := range l.Cfg.IncludePaths {
						var innerAbsPath = filepath.Join(includePath, path)

						if _, err := os.Stat(innerAbsPath); errors.Is(err, os.ErrNotExist) {
							continue
						}

						fmt.Println("Found included path ", path, "at absolute location ", innerAbsPath)
						absPath = innerAbsPath
						break
					}
					if len(absPath) == 0 {
						fmt.Println("Unable to find included path: ", path)
						continue
					}

					result.Imports = append(result.Imports, language.ImportEntry{
						// Symbols: []string{path},
						AbsPath: absPath,
						All:     true,
					})
				}
		*/
	}

	return &result, nil
}

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	fmt.Println("Banchod export ", file)
	var result language.ExportsResult

	for _, statement := range file.Content.([]Statement) {
		var path string = ""
		if statement.Quoted != nil {
			path = statement.Quoted.IncludedFile
		} else if statement.Angled != nil {
			path = statement.Angled.IncludedFile
		}
		if len(path) == 0 {
			continue
		}

		result.Exports = append(result.Exports, language.ExportEntry{
			Symbols: []language.ExportSymbol{{Original: filepath.Base(path)}}, AbsPath: file.AbsPath,
		})
		/*
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
		*/
	}

	return &result, nil
}

var Extensions = []string{"h", "cpp", "cppm", "ixx"}
