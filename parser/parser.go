package parser

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/lexer"
)

type Parser struct {
    tokens []lexer.Token
    pos    int
}

func New(tokens []lexer.Token) *Parser {
    return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) cur() lexer.Token {
    if p.pos >= len(p.tokens) {
        return lexer.Token{Type: lexer.EOF, Lit: ""}
    }
    return p.tokens[p.pos]
}

func (p *Parser) next() lexer.Token {
    p.pos++
    return p.cur()
}

func (p *Parser) expect(tt lexer.TokenType) (lexer.Token, error) {
    tok := p.cur()
    if tok.Type != tt {
        return tok, fmt.Errorf("expected %s, got %s (%s)", tt, tok.Type, tok.Lit)
    }
    p.next()
    return tok, nil
}

func (p *Parser) ParseProgram() (*Program, error) {
    prog := &Program{}
    for p.cur().Type != lexer.EOF {
        if p.cur().Type == lexer.IMPORT {
            importPath, err := p.parseImport()
            if err != nil {
                return nil, err
            }
            prog.Imports = append(prog.Imports, importPath)
            continue
        }
        if p.cur().Type == lexer.FUNCTION || p.cur().Type == lexer.IDENT && p.cur().Lit == "fn" {
            fn, err := p.parseFunction()
            if err != nil {
                return nil, err
            }
            prog.Functions = append(prog.Functions, fn)
            continue
        }
        // support 'fn' using the keyword map
        if p.cur().Type == lexer.FUNCTION {
            fn, err := p.parseFunction()
            if err != nil {
                return nil, err
            }
            prog.Functions = append(prog.Functions, fn)
            continue
        }
        // fallback: try to parse function if 'fn' as IDENT
        if p.cur().Type == lexer.IDENT && p.cur().Lit == "fn" {
            // convert to FUNCTION
            p.tokens[p.pos].Type = lexer.FUNCTION
            continue
        }
        // skip unknown/illegal tokens
        p.next()
    }
    return prog, nil
}

func (p *Parser) parseFunction() (*Function, error) {
    // expect 'fn'
    if p.cur().Type == lexer.IDENT && p.cur().Lit == "fn" {
        p.tokens[p.pos].Type = lexer.FUNCTION
    }
    if _, err := p.expect(lexer.FUNCTION); err != nil {
        return nil, err
    }
    nameTok, err := p.expect(lexer.IDENT)
    if err != nil {
        return nil, err
    }
    // expect ()
    if _, err := p.expect(lexer.LPAREN); err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.RPAREN); err != nil {
        return nil, err
    }
    // optional return type: -> IDENT
    retType := "int"
    if p.cur().Type == lexer.ARROW {
        p.next()
        if p.cur().Type == lexer.IDENT {
            retType = p.cur().Lit
            p.next()
        } else {
            return nil, fmt.Errorf("expected return type identifier after '->'")
        }
    }
    // body
    if _, err := p.expect(lexer.LBRACE); err != nil {
        return nil, err
    }
    body := &BlockStatement{}
    for p.cur().Type != lexer.RBRACE && p.cur().Type != lexer.EOF {
        stmt, err := p.parseStatement()
        if err != nil {
            return nil, err
        }
        if stmt != nil {
            body.Statements = append(body.Statements, stmt)
        }
    }
    if _, err := p.expect(lexer.RBRACE); err != nil {
        return nil, err
    }
    return &Function{Name: nameTok.Lit, ReturnType: retType, Body: body}, nil
}

func (p *Parser) parseStatement() (Statement, error) {
    switch p.cur().Type {
    case lexer.RETURN:
        return p.parseReturn()
    case lexer.VAR:
        return p.parseVar()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseReturn() (Statement, error) {
    p.next()
    expr, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    // optional semicolon
    if p.cur().Type == lexer.SEMICOLON {
        p.next()
    }
    return &ReturnStatement{Value: expr}, nil
}

func (p *Parser) parseVar() (Statement, error) {
    p.next()
    nameTok, err := p.expect(lexer.IDENT)
    if err != nil {
        return nil, err
    }
    varType := "int"
    if p.cur().Type == lexer.COLON {
        p.next()
        if p.cur().Type == lexer.IDENT {
            varType = p.cur().Lit
            p.next()
        } else {
            return nil, fmt.Errorf("expected type identifier after ':'")
        }
    }
    if _, err := p.expect(lexer.ASSIGN); err != nil {
        return nil, err
    }
    expr, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    if p.cur().Type == lexer.SEMICOLON {
        p.next()
    }
    return &VarStatement{Name: nameTok.Lit, Type: varType, Value: expr}, nil
}

func (p *Parser) parseExpressionStatement() (Statement, error) {
    expr, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    if p.cur().Type == lexer.SEMICOLON {
        p.next()
    }
    return &ExpressionStatement{Expr: expr}, nil
}

func (p *Parser) parseExpression() (Expression, error) {
    return p.parseExpressionWithPrecedence(0)
}

var precedences = map[lexer.TokenType]int{
    lexer.EQ:       1,
    lexer.NOT_EQ:   1,
    lexer.LT:       2,
    lexer.GT:       2,
    lexer.PLUS:     3,
    lexer.MINUS:    3,
    lexer.SLASH:    4,
    lexer.ASTERISK: 4,
}

func (p *Parser) parseExpressionWithPrecedence(precedence int) (Expression, error) {
    tok := p.cur()
    var left Expression
    switch tok.Type {
    case lexer.IDENT:
        ident := &Identifier{Value: tok.Lit}
        p.next()
        if p.cur().Type == lexer.LPAREN {
            p.next()
            var args []Expression
            for p.cur().Type != lexer.RPAREN && p.cur().Type != lexer.EOF {
                arg, err := p.parseExpression()
                if err != nil {
                    return nil, err
                }
                args = append(args, arg)
                if p.cur().Type == lexer.COMMA {
                    p.next()
                }
            }
            if _, err := p.expect(lexer.RPAREN); err != nil {
                return nil, err
            }
            left = &CallExpression{Function: ident, Args: args}
        } else if p.cur().Type == lexer.LBRACKET {
            // convert slice/index syntax `id[expr]` or `id[start,end]` into a
            // call to runtime helper `Slice(id, start, end)` where end is -1 when omitted.
            p.next()
            // parse start expression
            startExpr, err := p.parseExpression()
            if err != nil {
                return nil, err
            }
            var endExpr Expression
            if p.cur().Type == lexer.COMMA {
                p.next()
                e, err := p.parseExpression()
                if err != nil {
                    return nil, err
                }
                endExpr = e
            } else {
                // end omitted -> use integer literal -1 to indicate 'to end'
                endExpr = &IntegerLiteral{Value: "-1"}
            }
            if _, err := p.expect(lexer.RBRACKET); err != nil {
                return nil, err
            }
            // build call: Slice(ident, startExpr, endExpr)
            args := []Expression{ident, startExpr, endExpr}
            left = &CallExpression{Function: &Identifier{Value: "Slice"}, Args: args}
        } else {
            left = ident
        }
    case lexer.INT:
        lit := &IntegerLiteral{Value: tok.Lit}
        p.next()
        left = lit
    case lexer.STRING:
        lit := &StringLiteral{Value: tok.Lit}
        p.next()
        left = lit
    case lexer.LPAREN:
        p.next()
        expr, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        if _, err := p.expect(lexer.RPAREN); err != nil {
            return nil, err
        }
        left = expr
    default:
        return nil, fmt.Errorf("unexpected token in expression: %s (%s)", tok.Type, tok.Lit)
    }

    for p.cur().Type != lexer.SEMICOLON && p.cur().Type != lexer.COMMA && p.cur().Type != lexer.RPAREN && p.cur().Type != lexer.EOF {
        curPrec, ok := precedences[p.cur().Type]
        if !ok || curPrec <= precedence {
            break
        }
        opTok := p.cur()
        p.next()
        right, err := p.parseExpressionWithPrecedence(curPrec)
        if err != nil {
            return nil, err
        }
        left = &InfixExpression{Left: left, Operator: opTok.Lit, Right: right}
    }
    return left, nil
}

// parseImport parses an import statement: import "path/to/file"
func (p *Parser) parseImport() (string, error) {
    // expect 'import' keyword
    if _, err := p.expect(lexer.IMPORT); err != nil {
        return "", err
    }
    
    // expect string literal with import path
    if p.cur().Type != lexer.STRING {
        return "", fmt.Errorf("expected string literal after 'import', got %s", p.cur().Type)
    }
    
    importPath := p.cur().Lit
    p.next()
    
    // expect semicolon
    if _, err := p.expect(lexer.SEMICOLON); err != nil {
        return "", err
    }
    
    return importPath, nil
}
