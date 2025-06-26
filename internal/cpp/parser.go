package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

/*
type ImportStatement struct {
	Symbols []string `"#include" "@Ident"`
	From    string   `"from" @(Ident|Punctuation|"/")*`
}
*/

/*
type ExportStatement struct {
	Symbol string `"export" @Ident`
}

type Statement struct {
	Import *ImportStatement `@@ |`
	Export *ExportStatement `@@`
}
*/

type QuotedInclude struct {
	IncludeToken string `@Include`
	IncludedFile string `@String`
}

type AngledInclude struct {
	IncludeToken string `@Include`
	IncludedFile string `@Angled`
}

type Statement struct {
	Quoted *QuotedInclude `@@`
	Angled *AngledInclude `|@@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"Include", "#include"},
			{"String", `"[^"]+"`},
			{"Angled", `<[^>]+>`},
			{"Ident", `[A-Za-z_][A-Za-z0-9_]*`},
			{"Whitespace", `\s+`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Unquote("String", "Angled"),
		participle.Elide("Whitespace"),
	)
)
