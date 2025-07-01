package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type UsingDirective struct {
	Namespace string `"using" "namespace" @Ident ("::" @Ident)* ";"`
}

type IncludeDirective struct {
	Quoted *string `"#" "include" (@String`
	Angled *string `| @AngledInclude)`
}

type Statement struct {
	Include             *IncludeDirective `@@`
	IgnoredPreprocessor *string           `| @PreprocessorLine`
	Using               *UsingDirective   `| @@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": {
			lexer.Include("Comment"),

			{"PreprocessorLine", `#[^\r\n]*`, nil},

			{"Using", `using\b`, nil},
			{"Namespace", `namespace\b`, nil},
			{"Include", `include\b`, nil},

			{"String", `"(?:[^"\\]|\\.)*"`, nil},
			{"AngledInclude", `<[^>\r\n]*>`, nil},

			{"NamespaceAccess", `::`, nil},
			{"Semi", `;`, nil},
			{"Hash", `#`, nil},

			{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`, nil},

			{"Whitespace", `[ \t]+`, nil},
			{"Newline", `[\r\n]+`, nil},
		},

		"Comment": {
			{"LineComment", `//[^\r\n]*`, nil},
			{"BlockComment", `\s*/\*[\s\S\n\r]*?\*/`, nil},
		},
	})

	parser = participle.MustBuild[File](
		participle.Lexer(def),
		participle.Unquote("String"),
		participle.Elide("Whitespace", "Newline", "LineComment", "BlockComment"),
	)
)
