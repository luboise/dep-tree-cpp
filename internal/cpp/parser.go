package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type NamespaceDef struct {
	Name       *TypeName   `"namespace" @@`
	Statements []Statement `"{" "\n"? @@* "}"`
}

type ClassFwd struct {
	Name string `"class" @Ident ";"`
}

type UsingDirective struct {
	Namespace string `"namespace" @Ident ("::" @Ident)* ";"`
}

// https://en.cppreference.com/w/cpp/language/namespace.html#Using-declarations
// eg.
// using std::vector, std::string, mynamespace::foo, mynamespace::bar;
type UsingDeclaration struct {
	// Tokens []string `@Ident ("," @Ident)*`
	Tokens []TypeName `@@ ("," @@)* ";"`
}

type TypeName struct {
	Name string `@NamespaceAccess? @Ident (@NamespaceAccess @Ident)*`
}

type TypeAlias struct {
	Alias string   `@Ident "="`
	Type  TypeName `@@ ";"`
}

type IncludeDirective struct {
	Quoted *string `"#" "include" (@String`
	Angled *string `| @AngledInclude)`
}

type UsingStatement struct {
	UsingDirective   *UsingDirective   `"using" (@@`
	UsingDeclaration *UsingDeclaration `| @@`
	TypeAlias        *TypeAlias        `| @@)`
}

type Statement struct {
	Include             *IncludeDirective `@@`
	IgnoredPreprocessor *string           `| @PreprocessorLine`
	Namespace           *NamespaceDef     `| @@`
	Using               *UsingStatement   `| @@`
	ClassDef            *ClassFwd         `| @@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": {
			lexer.Include("Comment"),

			{"PreprocessorLine", `#[^\r\n]*`, nil},

			{"Class", `class\b`, nil},

			{"Using", `using\b`, nil},
			{"Namespace", `namespace\b`, nil},
			{"Include", `include\b`, nil},

			{"String", `"(?:[^"\\]|\\.)*"`, nil},
			{"AngledInclude", `<[^>\r\n]*>`, nil},

			{"NamespaceAccess", `::`, nil},
			{"Equals", `=`, nil},
			{"Semi", `;`, nil},
			{"Comma", `,`, nil},
			{"Hash", `#`, nil},
			{"AngledL", `<`, nil},
			{"AngledR", `>`, nil},

			{"CurlyOpen", `{`, nil},
			{"CurlyClose", `}`, nil},

			{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`, nil},

			{"Whitespace", `[ \t]+`, nil},
			{"Newline", `[\r\n]+`, nil},
			{"Anything", ".*\n?", nil},
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
