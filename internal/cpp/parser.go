package cpp

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

/*
type ImportStatement struct {
	Symbols []string `"import" @Ident ("," @Ident)*`
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
	IncludedFile string `"#include" @String`
}

type AngledInclude struct {
	IncludedFile string `"#include" @Angled`
}

/*
type QuotedInclude struct {
	IncludeToken string `@Include`
	IncludedFile string `@String`
}

type AngledInclude struct {
	IncludeToken string `@Include`
	IncludedFile string `@Angled`
}
*/

type LineComment struct {
	// Comment string `@LineCommentStart @LineCommentEnd`
	Comment string `@LineComment`
}

type BlockComment struct {
	Comment string `@BlockComment`
}

type Comment struct {
	LineComment  *LineComment  `@@`
	BlockComment *BlockComment `| @@`
}

type Statement struct {
	Quoted  *QuotedInclude `@@`
	Angled  *AngledInclude `| @@`
	Comment *Comment       `| @@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"BlockComment", `/\*(.|\n)+\*/`},
			{"Include", "#include"},
			{"KeyWord", "(export|import|from)"},
			{"Punctuation", `[,\.]`},
			{"Ident", `[a-zA-Z]+`},
			{"Newline", `\n+`},
			{"String", `"[^"]+"`},
			{"Angled", `<[^>]+>`},
			{"LineComment", `//[^\r\n]*`},
			{"Whitespace", `[ \t]+`},

			// {"LineCommentStart", `//`},
			// {"LineCommentEnd", `.+$`},

			/*
				{"Include", "#include"},
				{"KewWord", "(int|float)"},
				{"Ident", `[A-Za-z_][A-Za-z0-9_]*`},
				{"Whitespace", `[ \t]+`},
			*/
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Unquote("String", "Angled"),
		participle.Elide("Newline", "Whitespace"),
	)
)
