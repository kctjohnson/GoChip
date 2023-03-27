package parser

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	// Sets the char to the first character, position to 0, and read position to 1
	l.readChar()
	return l
}

func (l *Lexer) Rewind() {
	l.position = 0
	l.readPosition = 0
	l.ch = 0
	l.readChar()
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		tok = NewToken(LBRACKET, l.ch)
	case ']':
		tok = NewToken(RBRACKET, l.ch)
	case ',':
		tok = NewToken(COMMA, l.ch)
	case ':':
		tok = NewToken(COLON, l.ch)
	case '0':
		tok = l.readValue()
		return tok
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LoopupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok = l.readValue()
			return tok
		} else {
			tok = NewToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Reads the full value, determining if it's hex or not
func (l *Lexer) readValue() Token {
	position := l.position
	token := Token{}
	token.Type = DECIMAL
	for isDigit(l.ch) || l.ch == 'x' {
		if l.ch == 'x' {
			token.Type = HEX
		}
		l.readChar()
	}
	token.Literal = l.input[position:l.position]
	return token
}

// Returns true if the current character is a digit
func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f')
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\r' && l.ch != '\n' {
		l.readChar()
	}
	return l.input[position:l.position]
}
