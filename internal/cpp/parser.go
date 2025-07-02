package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Expr struct {
	String *string `@String`
}

type NamespaceDef struct {
	Name       *TypeName   `"namespace" @@`
	Statements []Statement `"{" "\n"? @@* "}"`
}

type ClassFwd struct {
	Name string `"class" @Ident ";"`
}

type FunctionParameter struct {
	Type  *QualifiedTypeName `@@`
	Name  *string            `@Ident?`
	Value *Expr              `(("=" @@) | ("{" @@ "}"))?`
}

type FunctionDeclaration struct {
	Qualifiers []*string `("[" "[" ("nodiscard"|"deprecated") ("(" @String ")")? "]" "]")*`
	// Qualifiers        []*string            `(("[[" @FunctionQualifier ("(" @String ")") "]]")|@FunctionQualifier)*`
	LeadingReturnType TypeName             `@@`
	Name              string               `@Ident`
	Parameters        []*FunctionParameter `"(" (@@ ("," @@)*)? ")" ";"`
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

type TemplatedType struct {
	Types []TypeName `"<" @@ ("," @@)* ">"`
}

type TypeAlias struct {
	Alias string   `@Ident "="`
	Type  TypeName `@@ ";"`
}

type TypeExtensionFragment struct {
	NamespaceAccess *string     `("::" @Ident)`
	TemplatedType   []*TypeName `| ("<" @@ (","@@)* ">")`
}

type QualifiedTypeName struct {
	IsConst         bool      `"const"?`
	IsVolatile      bool      `"volatile"?`
	Type            TypeName  `@@`
	CompoundSymbols []*string `(@Ampersand|@Asterisk)*`
}

type TypeName struct {
	GlobalNamespace bool                     `"::"?`
	Name            string                   `@Ident`
	Fragments       []*TypeExtensionFragment `@@*`
}

type IncludeDirective struct {
	Quoted *string `"#" "include" (@String`
	Angled *string `| (@AngledL @Ident @AngledR))`
}

type UsingStatement struct {
	UsingDirective   *UsingDirective   `"using" (@@`
	UsingDeclaration *UsingDeclaration `| @@`
	TypeAlias        *TypeAlias        `| @@)`
}

type Statement struct {
	Include             *IncludeDirective    `@@`
	IgnoredPreprocessor *string              `| @PreprocessorLine`
	Namespace           *NamespaceDef        `| @@`
	Using               *UsingStatement      `| @@`
	ClassDef            *ClassFwd            `| @@`
	EmptyStatement      *string              `| @Semi`
	FunctionDeclaration *FunctionDeclaration `| @@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": {
			lexer.Include("Comment"),

			{"FunctionQualifier", `"nodiscard"|"deprecated"`, nil},

			{"PreprocessorLine", `#[^\r\n]*`, nil},

			{"Class", `class\b`, nil},

			{"Using", `using\b`, nil},
			{"Namespace", `namespace\b`, nil},
			{"Include", `include\b`, nil},

			{"String", `"(?:[^"\\]|\\.)*"`, nil},
			// {"AngledInclude", `<[^>\r\n]*>`, nil},

			{"NamespaceAccess", `::`, nil},
			{"Equals", `=`, nil},
			{"Semi", `;`, nil},
			{"Comma", `,`, nil},
			{"Hash", `#`, nil},
			{"AngledL", `<`, nil},
			{"AngledR", `>`, nil},

			{"Asterisk", `\*`, nil},

			{"RValue", `&&`, nil},
			{"Ampersand", `&`, nil},

			{"CurlyL", `\{`, nil},
			{"CurlyR", `\}`, nil},

			{"SquareL", `\[`, nil},
			{"SquareR", `\]`, nil},

			{"RoundL", `\(`, nil},
			{"RoundR", `\)`, nil},

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
