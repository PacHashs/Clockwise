package parser

import "codeberg.org/clockwise-lang/clockwise/lexer"

// IfStatement represents a simple if/else statement.
type IfStatement struct {
    Condition Expression
    Consequent *BlockStatement
    Alternative *BlockStatement
}

// WhileStatement represents a while loop.
type WhileStatement struct {
    Condition Expression
    Body *BlockStatement
}

// parse helpers: these are small API additions; parser will call them
// if extended control parsing is desired.
func (p *Parser) parseIf() (Statement, error) {
    // expects current token is IF
    if _, err := p.expect(lexer.IF); err != nil {
        return nil, err
    }
    // expect '(' condition ')'
    if _, err := p.expect(lexer.LPAREN); err != nil {
        return nil, err
    }
    cond, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.RPAREN); err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.LBRACE); err != nil {
        return nil, err
    }
    cons := &BlockStatement{}
    for p.cur().Type != lexer.RBRACE && p.cur().Type != lexer.EOF {
        st, err := p.parseStatement()
        if err != nil {
            return nil, err
        }
        if st != nil {
            cons.Statements = append(cons.Statements, st)
        }
    }
    if _, err := p.expect(lexer.RBRACE); err != nil {
        return nil, err
    }
    var alt *BlockStatement
    if p.cur().Type == lexer.ELSE {
        p.next()
        if _, err := p.expect(lexer.LBRACE); err != nil {
            return nil, err
        }
        alt = &BlockStatement{}
        for p.cur().Type != lexer.RBRACE && p.cur().Type != lexer.EOF {
            st, err := p.parseStatement()
            if err != nil {
                return nil, err
            }
            if st != nil {
                alt.Statements = append(alt.Statements, st)
            }
        }
        if _, err := p.expect(lexer.RBRACE); err != nil {
            return nil, err
        }
    }
    return &IfStatement{Condition: cond, Consequent: cons, Alternative: alt}, nil
}

func (p *Parser) parseWhile() (Statement, error) {
    if _, err := p.expect(lexer.WHILE); err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.LPAREN); err != nil {
        return nil, err
    }
    cond, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.RPAREN); err != nil {
        return nil, err
    }
    if _, err := p.expect(lexer.LBRACE); err != nil {
        return nil, err
    }
    body := &BlockStatement{}
    for p.cur().Type != lexer.RBRACE && p.cur().Type != lexer.EOF {
        st, err := p.parseStatement()
        if err != nil {
            return nil, err
        }
        if st != nil {
            body.Statements = append(body.Statements, st)
        }
    }
    if _, err := p.expect(lexer.RBRACE); err != nil {
        return nil, err
    }
    return &WhileStatement{Condition: cond, Body: body}, nil
}
