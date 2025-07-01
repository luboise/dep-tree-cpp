package cpp

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

/*

type Declaration struct {
	Namespace *NamespaceDef   "@@"
	TypeAlias *TypeAlias      "| @@"
	Using     *UsingStatement "| @@"
	Fwd       *FwdDec         "| @@"
	Function  *FnDec          "| @@"
}

type Type struct {
	IsConst bool   `"const"?`
	Name    string `("::"? @Ident | "<" @Ident ">")+ ("&"+ | "*"+)?`
}

type Exp struct {
	Left  string `@Ident`
	Right ExpR   `@@`
}

type ExpR struct {
	RVal *Value `"=" @@`
}

type String struct {
	String string `@String`
}

type Value struct {
	Function *FunctionCall `@@`
	String   *String       `| @@`
}

type FunctionCall struct {
	Name      string   `@Ident(`
	RoundArgs *ArgList `( "(" @@ ")" ) |`
	CurlyArgs *ArgList `( "{" @@ "}" )  )`
}

type Arg struct {
}

type ArgList struct {
	Values []Value `((@@ ",")* (@@)?)`
}

type IdentifierList struct {
	Identifiers []string `@Ident ("," @Ident)*`
}

type FwdDec struct {
	Class *ClassFwd `@@`
}

type FunctionType struct {
	Type Type    `@@`
	Name *string `@Ident?`

	// Value   *string `("=" @String)?`
	Value *Value `("=" @@)?`
}

type FnParam struct {
	Type  string `"const"? @Ident`
	Name  string `@Ident [\*&]?`
	Value string `"=" @Ident`
}

type FnPreSpecifiers struct {
	Specifiers []string `@Ident*`
}

type FnDec struct {
	PreSpecifiers      *FnPreSpecifiers `@@?`
	LeadingReturnType  Type             `@@`
	Name               string           `@Ident "("`
	Parameters         []FunctionType   `( @@ ( "," @@ )* )? ")"`
	PostSpecifiers     []string         `@Ident*`
	TrailingReturnType *string          `("->" @Ident)? ";"`
}

type ClassFwd struct {
	Name string `"class" @Ident ";"`
}

type VariableDeclaration struct {
	Type      string `@Ident`
	Name      string `@Ident`
	Semicolon string `@Semi`
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
	// Tokens []string `@Ident ("," @Ident)*`
	Tokens *IdentifierList `@@`
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

type Declaration struct {
	VarDec   *VariableDeclaration "@@"
	ClassDec *ClassDeclaration    "| @@"
}

type NamespaceDef struct {
	Name  string        `"namespace" @Ident "{"`
	Items []Declaration `@@* "}"`
}

type Include struct {
	Angled *string `"#include" @Angled @LineComment?`
	Quoted *string `| "#include" @String @LineComment?`
}
*/

type QuotedInclude struct {
	IncludedFile string `@QuotedInclude`
}

type AngledInclude struct {
	IncludedFile string `@AngledInclude`
}

type Statement struct {
	Quoted *QuotedInclude `@@`
	Angled *AngledInclude `| @@`
	// Empty   bool     `| (@Semi|"\n")` // Accept empty statements
}

type File struct {
	Statements []Statement `@@*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"QuotedInclude", `#include\s+"[^"]+"`},
			{"AngledInclude", `#include\s+<[^<]+>`},

			// {"BadPreprocessor", "^#([^i]|i[^n]|in[^c]|inc[^l]|incl[^u]|inclu[^d]|includ[^e])"},
			// {"Pragma", "#pragma.*\n"},
			// {"Define", "#define.*\n"},
			// {"BadLine", "^[^#].*\n"},
			// {"EmptyLine", "^\n$"},
			// {"Brother", "[^\n]+"},
			// {"KeyWord", "(const|export|import|using|namespace|class|struct)"},
			// {"Punctuation", `[,\.\{\}=\(\)&]`},
			// {"Semi", `;`},
			// {"Ident", `[a-zA-Z]+`},
			// {"Newline", `[\n\r]+`},
			// {"Ident", `([_a-zA-Z][a-zA-Z0-9]*::)*[_a-zA-Z0-9]+`},
			{"LineComment", `//[^\r\n]*`},
			{"BlockComment", `/\*(.|\n)+\*/`},
			{"Whitespace", `\s+`},
			{"Other", `.+`},

			// Elided rules
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		// participle.Unquote("String", "Angled"),
		participle.Elide("LineComment", "BlockComment", "Whitespace", "Other"),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			token.Value = strings.Replace(token.Value, "#include", "", -1)
			token.Value = strings.Replace(token.Value, `"`, "", -1)
			token.Value = strings.Replace(token.Value, `<`, "", -1)
			token.Value = strings.Replace(token.Value, `>`, "", -1)
			token.Value = strings.Replace(token.Value, ` `, "", -1)
			return token, nil
		}, "QuotedInclude", "AngledInclude"),
	)
)
