package parse

import (
	"fmt"
	"strings"
)

type Parser struct {
	lexer  *Lexer
	tokens []Token
}

func NewParser(input string) *Parser {
	lowerInput := strings.ToLower(input)
	return &Parser{
		lexer: NewLexer(lowerInput),
	}
}

func (p *Parser) ReadTokens() {
	p.tokens = []Token{}
	for {
		tok := p.lexer.NextToken()
		if tok.Type == EOF {
			p.tokens = append(p.tokens, tok)
			break
		} else if tok.Type == ILLEGAL {
			fmt.Printf("Illegal Token: %s\n", tok.Literal)
			p.tokens = append(p.tokens, tok)
			break
		}
		p.tokens = append(p.tokens, tok)
	}
}

func (p *Parser) Rewind() {
	p.lexer.Rewind()
}

func (p Parser) PrintTokens() {
	for _, tok := range p.tokens {
		fmt.Printf("%#v\n", tok)
	}
}
