package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	lexer  *Lexer
	tokens []Token
	labels []string
}

func NewParser(input string) *Parser {
	lowerInput := strings.ToLower(input)
	return &Parser{
		lexer: NewLexer(lowerInput),
	}
}

func (p *Parser) ReadTokens() {
	p.labels = []string{}
	p.tokens = []Token{}

	// First pass gathers the tokens
	for {
		tok := p.lexer.NextToken()
		if tok.Type == EOF {
			p.tokens = append(p.tokens, tok)
			break
		} else if tok.Type == ILLEGAL {
			fmt.Printf("Illegal Token: %s\n", tok.Literal)
			p.tokens = append(p.tokens, tok)
			break
		} else if tok.Type == COLON && p.tokens[len(p.tokens)-1].Type == UNKNOWNIDENT {
			prevTokIndex := len(p.tokens) - 1
			p.tokens[prevTokIndex].Type = LABEL_DEF
			p.labels = append(p.labels, p.tokens[prevTokIndex].Literal)
		}
		p.tokens = append(p.tokens, tok)
	}

	// Second pass checks for undefined label references
	for i, tok := range p.tokens {
		if tok.Type == UNKNOWNIDENT {
			for _, label := range p.labels {
				if tok.Literal == label {
					p.tokens[i].Type = LABEL_REF
					break
				}
			}
		}
	}
}

func (p Parser) GetTokens() []Token {
	return p.tokens
}

func (p *Parser) Rewind() {
	p.lexer.Rewind()
}

func (p Parser) PrintTokens() {
	for _, tok := range p.tokens {
		fmt.Printf("%#v\n", tok)
	}
}
