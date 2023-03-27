package parser

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	UNKNOWNIDENT = "UNKNOWNIDENT"
	LBRACKET     = "LBRACKET"
	RBRACKET     = "RBRACKET"
	COMMA        = "COMMA"
	COMMENT      = "COMMENT"

	// Values
	HEX     = "HEX"
	DECIMAL = "DECIMAL"

	// Keywords
	ADP       = "ADP"
	DELAY     = "DELAY"
	SND_DELAY = "SND_DELAY"
	REG       = "REG"
	CLS       = "CLS"
	SYSCALL   = "SYSCALL"
	CALL      = "CALL"
	RET       = "RET"
	JMP       = "JMP"
	RJMP      = "RJMP"
	SEQ       = "SEQ"
	SNEQ      = "SNEQ"
	JKP       = "JKP"
	JKNP      = "JKNP"
	WK        = "WK"
	MOV       = "MOV"
	ADD       = "ADD"
	SUB       = "SUB"
	OR        = "OR"
	AND       = "AND"
	XOR       = "XOR"
	SHR       = "SHR"
	SHL       = "SHL"
	BRND      = "BRND"
	DRW       = "DRW"
	FX29      = "FX29"
	FX33      = "FX33"
	FX55      = "FX55"
	FX65      = "FX65"
)

var keywords = map[string]TokenType{
	"adp":       ADP,
	"delay":     DELAY,
	"snd_delay": SND_DELAY,
	"reg":       REG,
	"cls":       CLS,
	"syscall":   SYSCALL,
	"call":      CALL,
	"ret":       RET,
	"jmp":       JMP,
	"rjmp":      RJMP,
	"seq":       SEQ,
	"sneq":      SNEQ,
	"jkp":       JKP,
	"jknp":      JKNP,
	"wk":        WK,
	"mov":       MOV,
	"add":       ADD,
	"sub":       SUB,
	"or":        OR,
	"and":       AND,
	"xor":       XOR,
	"shr":       SHR,
	"shl":       SHL,
	"brnd":      BRND,
	"drw":       DRW,
	"fx29":      FX29,
	"fx33":      FX33,
	"fx55":      FX55,
	"fx65":      FX65,
}

func LoopupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return UNKNOWNIDENT
}
