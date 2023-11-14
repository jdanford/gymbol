package gymbol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type Reader struct {
	bufio.Reader
	vm        *VM
	line, col int
}

type Error struct {
	error
	line, col int
}

const READER_BUF_SIZE = 64

func Read(vm *VM, input io.Reader) (Value, error) {
	bufReader := bufio.NewReaderSize(input, READER_BUF_SIZE)
	reader := Reader{*bufReader, vm, 0, 0}
	value, err := reader.Value()

	if errors.Is(err, io.EOF) {
		return nil, reader.Error(errors.New("can't parse empty expression"))
	}

	if err != nil {
		return nil, err
	}

	extra, err := reader.Peek()

	if err != nil {
		if errors.Is(err, io.EOF) {
			return value, nil
		}

		return nil, err
	}

	return nil, reader.Error(fmt.Errorf("extra character: `%c`", extra))
}

func (reader *Reader) Error(err error) Error {
	return Error{err, reader.line, reader.col}
}

func (reader *Reader) Peek() (rune, error) {
	r, _, err := reader.Reader.ReadRune()
	// fmt.Printf("peek: %c\n", r)

	if err == nil {
		err = reader.UnreadRune()
	}

	return r, err
}

func (reader *Reader) Read() (rune, error) {
	r, n, err := reader.Reader.ReadRune()
	// fmt.Printf("read: %c\n", r)

	if err == nil {
		if r == '\n' {
			reader.line += 1
			reader.col = 0
		} else {
			reader.col += n
		}
	}

	return r, err
}

func (reader *Reader) Skip() {
	_, err := reader.Read()
	if err != nil {
		panic(err)
	}
}

func (reader *Reader) ReadWhile(test func(rune) bool) (string, error) {
	buf := new(bytes.Buffer)

	for {
		r, err := reader.Peek()
		if err != nil {
			return "", err
		}

		if !test(r) {
			break
		}

		r, err = reader.Read()
		if err != nil {
			return "", err
		}

		buf.WriteRune(r)
	}

	return buf.String(), nil
}

func (reader *Reader) ReadUntil(test func(rune) bool) (string, error) {
	negated := func(r rune) bool {
		return !test(r)
	}

	return reader.ReadWhile(negated)
}

func (reader *Reader) SkipWhitespace() error {
	_, err := reader.ReadWhile(unicode.IsSpace)
	return err
}

func (reader *Reader) Value() (Value, error) {
	// fmt.Printf("parsing value\n")
	reader.SkipWhitespace()
	next, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	var value Value
	if next == L_PAREN {
		value, err = reader.List()
	} else if IsQuote(next) {
		value, err = reader.Quoted()
	} else {
		value, err = reader.Atom()
	}

	if err != nil {
		return nil, err
	}

	reader.SkipWhitespace()
	return value, nil
}

func (reader *Reader) Atom() (Value, error) {
	// fmt.Printf("parsing atom\n")
	next, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	if unicode.IsDigit(next) || next == DASH {
		return reader.Number()
	}

	if next == D_QUOTE {
		return reader.String()
	}

	return reader.Symbol()
}

func (reader *Reader) Number() (Value, error) {
	// fmt.Printf("parsing number\n")
	next, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	negative := next == DASH
	if negative {
		reader.Skip()
	}

	s, err := reader.ReadWhile(unicode.IsDigit)
	if err != nil {
		return nil, err
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}

	if negative {
		f = -f
	}

	return Number(f), nil
}

func (reader *Reader) Symbol() (Value, error) {
	// fmt.Printf("parsing symbol\n")
	s, err := reader.ReadWhile(IsSymbol)
	if len(s) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return reader.vm.Symbols.Intern(s), nil
}

func (reader *Reader) String() (Value, error) {
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	for {
		chunk, err := reader.ReadUntil(func(r rune) bool { return r == D_QUOTE || r == BACKSLASH })
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteString(chunk)
		if err != nil {
			return nil, err
		}

		next, err := reader.Peek()

		if errors.Is(err, io.EOF) {
			s := buf.String()
			return nil, reader.Error(fmt.Errorf("unterminated string literal: `%s`", s))
		}

		if err != nil {
			return nil, err
		}

		if next == D_QUOTE {
			reader.Skip()
			break
		}

		if next == BACKSLASH {
			reader.Skip()
			next, err = reader.Peek()
			if err != nil {
				return nil, err
			}

			if next == 'r' {
				reader.Skip()
				next = '\r'
			} else if next == 't' {
				reader.Skip()
				next = '\t'
			} else if next == 'n' {
				reader.Skip()
				next = '\n'
			} else if next == 'f' {
				reader.Skip()
				next = '\f'
			} else if next == 'b' {
				reader.Skip()
				next = '\b'
			} else {
				next, _, err = reader.ReadRune()
				if err != nil {
					return nil, err
				}
			}

			_, err = buf.WriteRune(next)
			if err != nil {
				return nil, err
			}

			continue
		}

		next, _, err = reader.ReadRune()
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteRune(next)
		if err != nil {
			return nil, err
		}
	}

	s := buf.String()
	return String(s), nil
}

func (reader *Reader) Quoted() (Value, error) {
	quote, err := reader.Read()
	if err != nil {
		return nil, err
	}

	quoteType := quoteTypes[quote]
	next, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	if quoteType == UNQUOTE && next == AT {
		reader.Skip()
		quoteType = UNQUOTE_SPLICING
	}

	reader.SkipWhitespace() // TODO: remove this?

	quotedValue, err := reader.Value()
	if err != nil {
		return nil, err
	}

	if quotedValue == nil {
		quoteName := quoteNames[quoteType]
		return nil, reader.Error(fmt.Errorf("expected value after `%s`", quoteName))
	}

	return NewQuoted(quoteType, quotedValue), nil
}

func (reader *Reader) List() (Value, error) {
	// fmt.Printf("parsing list\n")
	next, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	if next != L_PAREN {
		return nil, reader.Error(fmt.Errorf("expected `%c`, got `%c`", L_PAREN, next))
	}

	reader.Skip()

	values := []Value{}
	for {
		value, err := reader.Value()
		if err != nil {
			return nil, err
		}

		if value == nil {
			break
		}

		values = append(values, value)
	}

	next, err = reader.Peek()
	if err != nil {
		return nil, err
	}

	if next != R_PAREN {
		return nil, reader.Error(fmt.Errorf("expected `%c`, got `%c`", R_PAREN, next))
	}

	reader.Skip()

	return reader.vm.Heap.AllocList(values)
}
