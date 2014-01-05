package markdown

import (
	"github.com/russross/blackfriday"
	. "github.com/yak-labs/chirp-lang"
)

func cmdMarkdown(fr *Frame, argv []T) T {
	input := []byte(Arg1(argv).String())
	return MkString(string(blackfriday.MarkdownCommon(input)))
}

func init() {
	if Safes == nil {
		Safes = make(map[string]Command, 333)
	}

	Safes["markdown"] = cmdMarkdown
}
