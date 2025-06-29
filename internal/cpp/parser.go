package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statement struct {
	Quoted *QuotedInclude `@@`
	Angled *AngledInclude `| @@`
	Dec    *Declaration   `| @@`
	Empty  bool           "| @Semi"
	// Ignored *Ignored        `| (@@ | Pragma) `
}

type QuotedInclude struct {
	IncludedFile string `"#include" @String`
}

type AngledInclude struct {
	IncludedFile string `"#include" @Angled`
}

type Declaration struct {
	Namespace *NamespaceDef   "@@"
	TypeAlias *TypeAlias      "| @@"
	Using     *UsingStatement "| @@"
	Fwd       *FwdDec         "| @@"
}

type FwdDec struct {
	Class *ClassFwd `@@`
}

type ClassFwd struct {
	Name string `"class" @Ident ";"`
}

type File struct {
	Statements []Statement `@@*`
}

type VariableDeclaration struct {
	Type      string `@Ident`
	Name      string `@Ident`
	Semicolon string `@Semi`
}

type UsingBruh struct {
	Tokens []string `"using" "namespace"? @Ident ("," "namespace"?@Ident)* ";"`
}

type UsingRValue struct {
	Value string `"=" @Ident`
}

// https://en.cppreference.com/w/cpp/language/namespace.html#Using-directives
// eg.
// using namespace std;
type UsingDirective struct {
	Namespace string `"namespace" @Ident`
}

// https://en.cppreference.com/w/cpp/language/namespace.html#Using-declarations
// eg.
// using std::vector, std::string, mynamespace::foo, mynamespace::bar;
type UsingDeclaration struct {
	Tokens []string `@Ident ("," @Ident)*`
}

type TypeAlias struct {
	Identifier string `@Ident "="`
	TypeID     string `@Ident @Angled?`
}

type UsingStatement struct {
	Alias       *TypeAlias        `"using" (@@`
	Declaration *UsingDeclaration "| @@"
	Directive   *UsingDirective   `| @@) ";"`
}

/*
type Declaration struct {
	VarDec   *VariableDeclaration "@@"
	ClassDec *ClassDeclaration    "| @@"
}
*/

type NamespaceDef struct {
	Name  string        `"namespace" @Ident "{"`
	Items []Declaration `@@* "}"`
}

/*
type Ignored struct {
	Alias *TypeAlias `@@`
	// Using    *UsingBruh `| @@`
	NSDef    *NamespaceDef `| @@`
	Declared *DeclaredItem `| @@`
}
*/

/*
type ClassItem {
	Var *VariaVariableDeclaration `@Ident @Ident`
}
*/

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"Include", "#include"},
			{"KeyWord", "(export|import|using|namespace|class|struct)"},
			{"Punctuation", `[,\.\{\}=]`},
			{"Semi", `;`},
			// {"Ident", `[a-zA-Z]+`},
			{"Newline", `\n+`},
			{"Ident", `([_a-zA-Z][a-zA-Z0-9]*::)*[_a-zA-Z0-9]+`},
			{"String", `"[^"]+"`},
			{"Angled", `<[^>]+>`},
			{"Whitespace", `[ \t]+`},

			// Elided rules
			{"LineComment", `//[^\r\n]*`},
			{"BlockComment", `/\*(.|\n)+\*/`},
			{"Pragma", "#pragma.*\n"},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Unquote("String", "Angled"),
		participle.Elide("LineComment", "BlockComment", "Newline", "Whitespace", "Pragma"),
	)
)
