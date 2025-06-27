package cpp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	simple_tests := []struct {
		Name       string
		Input      string
		Statements []Statement
	}{
		{
			Name:  "Quoted Include",
			Input: `#include "file.h"`,
			Statements: []Statement{{
				Quoted: &QuotedInclude{"file.h"},
			}},
		},
		{
			Name:  "Angled Include",
			Input: `#include <vector>`,
			Statements: []Statement{{
				Angled: &AngledInclude{"vector"},
			}},
		},
		{
			Name: "Namespace with using and forward declarations",
			Input: `
namespace Foo {

class ForwardedClass;
using BarPtr = Ptr<Bar>;

}  // namespace Presto
			`,
			Statements: []Statement{
				{
					Dec: &Declaration{
						Namespace: &NamespaceDef{
							Name: "Foo",
							Items: []Declaration{
								{
									Fwd: &FwdDec{
										Class: &ClassFwd{
											Name: "ForwardedClass",
										},
									},
								},
								{
									Using: &UsingStatement{
										Alias: &TypeAlias{
											Identifier: "BarPtr",
											TypeID:     "Ptr<Bar>",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	file_tests := []struct {
		Name     string
		File     string
		Expected []Statement
	}{
		{
			Name: "Multiple includes in one file",
			File: "multi_includes.h",
			Expected: []Statement{
				{
					Quoted: &QuotedInclude{"file.h"},
				},
				{
					Angled: &AngledInclude{"vector"},
				}},
		},
	}

	for _, tt := range simple_tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := parser.ParseBytes("", []byte(tt.Input))
			a.NoError(err)

			a.Equal(tt.Statements, result.Statements)
		})
	}

	const importsTestFolder = ".test_files"

	wd, _ := os.Getwd()
	for _, tt := range file_tests {
		t.Run(tt.Name, func(t *testing.T) {

			a := require.New(t)

			f := filepath.Join(wd, importsTestFolder, tt.File)

			bytes, err := os.ReadFile(f)
			a.NoError(err)

			result, err := parser.ParseBytes("", bytes)
			a.NoError(err)

			a.Equal(tt.Expected, result.Statements)
		})
	}
}
