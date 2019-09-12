package parse

import (
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log"
	"reflect"
	"strings"
)

const (
	errMalformedTableMsg         = "Parsing malformed table"
	errMalformedRowMsg           = "Parsing malformed row"
	errMalformedColMsg           = "Parsing malformed column"
	errMalformedEntity           = "Parsing encountered a malform state. Expected a non-empty stack."
	errUnexpectedTokenMsgFmt     = "Parsing an unexpected token. Token: '%v'"
	errMismatchedTokenMsgFmt     = "Parsing a mismatched token. Expected: '%v', but got: '%v'"
	errImpossibleCaseMsg         = "Parser encountered an impossible case"
	errUnexpectedTokenTypeMsgFmt = "Parsing encountered an unexpected token of type: %s"
	errUnknownTokenTypeMsgFmt    = "Parsing unexpected token of type: %s"
	errUnknownMsg                = "Unknown parsing error encountered"
)

type table []row
type row []string

// parserTokenizer is an interface for a regular html tokenizer created for testing
// purposes.
type parserTokenizer interface {
	Next() html.TokenType
	Token() html.Token
	Err() error
	Raw() []byte
}

// parserStack is an interface for a regular stack created for testing purposes.
type parserStack interface {
	Push(interface{})
	Pop() interface{}
	Len() int
	Peek() interface{}
}

// parserContext is a container for managing the data structures required
// for parsing nested entities.
type parserContext struct {
	sb strings.Builder
	st parserStack
	tk parserTokenizer
}

// tagCtr is a container for counting all token entities relevant to table parsing.
type tagCtr struct {
	col   int
	row   int
	table int
}

type persistedTokenizer struct {
	token     html.Token
	err       error
	tokenizer *html.Tokenizer
	raw       []byte
}

func (s *persistedTokenizer) Next() html.TokenType {
	result := s.tokenizer.Next()
	s.token = s.tokenizer.Token()
	s.err = s.tokenizer.Err()
	s.raw = s.tokenizer.Raw()
	return result
}

func (s *persistedTokenizer) Token() html.Token {
	return s.token
}

func (s *persistedTokenizer) Err() error {
	return s.err
}

func (s *persistedTokenizer) Raw() []byte {
	return s.raw
}

func newPersistedTokenizer(reader io.Reader) parserTokenizer {
	tk := html.NewTokenizer(reader)
	return &persistedTokenizer{
		tokenizer: tk,
	}
}

// parseDataTableHtmlTokenizer takes in a Reader and parses out strings
// by row and column into a gene-specific table and row struct type.
//
// It does not perform any validation of the table's data. For instance,
// whether a row is about headers or rows. Any nested tables get extracted
// as their own gene-specific table type, and appear before their parent
// table in the list of results.
//
// Any errors encountered in tokenizing the Reader output is returned. Since
// the tokenizer marks the End-Of-File (EOF) as an error token, this method does
// that check for the caller and returns back a nil error, if EOF is encountered.
func parseDataTableHtmlTokenizer(reader io.Reader) ([]table, error) {
	ctr := &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	ctx := &parserContext{
		sb: strings.Builder{},
		st: stack.New(),
		tk: newPersistedTokenizer(reader),
	}
	results := []table{}
	tableObj, rowObj := new(table), new(row)
	ctx.tk.Next()
	for {
		token := ctx.tk.Token()
		var err error
		switch token.Type {
		case html.ErrorToken:
			log.Printf("Handling Error Token: %+v\n", token)
			err = ctx.tk.Err()
			if err == io.EOF {
				return results, nil
			}

			if err != nil {
				return nil, err
			}
			return nil, errors.New(errUnknownMsg)
		case html.TextToken:
			log.Printf("Handling Text Tag: %+v\n", token)
			err = handleTextTagToken(ctx)
		case html.StartTagToken:
			log.Printf("Handling Starting Tag: %+v\n", token)
			handleStartTagToken(ctr, ctx)
		case html.EndTagToken:
			log.Printf("Handling End Tag: %+v\n", token)
			err = handleEndTagToken(rowObj, tableObj, &results, ctr, ctx)
		default:
			// Skipping all comments, self-closing tags e.g. <br />, and doctype tokens.
			log.Printf("Skipping Tag: %+v\n", token)
		}

		if err != nil {
			return nil, err
		}

		ctx.tk.Next()
	}
}

// handleTextTagToken extracts whatever text is at the current cursor
// position as long as the previous token in the parser stack is a
// column token (<td> or <th>). Otherwise, it skips the token completely,
// regardless of the token type.
//
// Returns an error if the current item on the stack is not a token type.
func handleTextTagToken(ctx *parserContext) error {
	// The only time a data entry is picked up, is if:
	// * Non-nested column
	// * Non-nested row
	// * Non-nested table
	//
	// Otherwise, the data is skipped.
	if ctx.st.Len() <= 0 {
		log.Println("Malformed Entity Encountered")
		return errors.New(errMalformedEntity)
	}

	token, tokenOk := ctx.st.Peek().(*html.Token)
	if !tokenOk {
		log.Println("Non html token discovered")
		return fmt.Errorf(errUnknownTokenTypeMsgFmt, reflect.TypeOf(token))
	}
	// Check that the parsing is in a column. The validity of the column is not checked here.
	if token.Type == html.StartTagToken && (token.DataAtom == atom.Td || token.DataAtom == atom.Th) {
		text := ctx.tk.Token().String()
		log.Printf("Extracting Text: %s\n", text)
		ctx.sb.WriteString(text)
	}

	return nil
}

// handleStartTagToken populates the token stack as with
// table, row, and/or column/header tokens only. Performs
// no validation on the stack. All validation occurs when
// the entity is closed out in handleEndTagToken.
func handleStartTagToken(ctr *tagCtr, ctx *parserContext) {
	token := ctx.tk.Token()
	if token.DataAtom == atom.Table {
		log.Println("Count Table")
		ctr.table++
	} else if token.DataAtom == atom.Tr {
		log.Println("Count Row")
		ctr.row++
	} else if token.DataAtom == atom.Td || token.DataAtom == atom.Th {
		log.Println("Count Column")
		ctr.col++
	} else {
		// Skip the push onto the stack. It's not an important token.
		log.Printf("Skipping Start Tag: %+v\n", token)
		return
	}

	log.Printf("Pushing Start Tag: %+v\n", token)
	ctx.st.Push(&token)
}

// handleEndTagToken is where conversion of the raw html into a gene-specific type happens.
// The token stack is popped, and the token at the tokenizer's cursor is compared and ensured
// to be the same. Should it be the same, either:
//  * Data is extracted.
//  * Gene-specific row is created.
//  * Gene-specific table is created and added to the results.
//
// Otherwise, an error is returned specific to one of the validation cases itemized below.
// Validation that's performed here:
//  * Tokens coming off the stack are of the HTML variety.
//  * Columns cannot be parsed outside of rows and tables.
//  * Rows cannot be parsed outside of tables.
//  * Stack should be empty iff. it's not inside a table.
//  * All tokens are of type table, row, and column/header.
//
// Errors are returned for each case above, should it fail.
func handleEndTagToken(rowObj *row, tableObj *table, results *[]table, ctr *tagCtr, ctx *parserContext) error {
	token := ctx.tk.Token()
	if token.DataAtom != atom.Table && token.DataAtom != atom.Tr && token.DataAtom != atom.Th && token.DataAtom != atom.Td {
		log.Printf("Skipping End Tag: %+v\n", token)
		return nil
	}

	tK, tKOk := ctx.st.Pop().(*html.Token)
	if !tKOk {
		err := fmt.Errorf(errUnexpectedTokenTypeMsgFmt, reflect.TypeOf(tK))
		log.Println(err)
		return err
	}

	if token.DataAtom == tK.DataAtom {
		log.Println("Valid Closing Tag")
		// For each of the tokens below, a match implies that structure is to be closed out and saved.
		// If the token enters one of the cases, either one of the cases must be true:
		//  * Table object was found and is closing out now.
		//  * Fragment of an object was found. If object counter is 0, process will
		//    error out.
		switch token.DataAtom {
		case atom.Table:
			log.Println("Reducing Table")
			ctr.table--
			excessTableTokens := ctr.table == 0 && ctx.st.Len() > 0
			lackingTableTokens := ctr.table > 0 && ctx.st.Len() == 0
			twilightZone := ctr.table < 0
			malformedTable := excessTableTokens || lackingTableTokens || twilightZone
			if malformedTable {
				log.Println(errMalformedTableMsg)
				return errors.New(errMalformedTableMsg)
			}

			if ctr.table == 0 {
				log.Println("Adding table to results")
				*results = append(*results, *tableObj)
				*tableObj = *new(table)
			}
		case atom.Tr:
			log.Println("Reducing Row")
			ctr.row--
			// There should never be the following cases:
			//  * Nested Rows: <tr><tr>...</tr></tr>
			//  * Standalone Rows. The row must be preceded by
			//    at least one element, which is expected to be
			//    <table>
			//
			//    No validation is performed on the last case at this
			//    stage. When those token's closing tags appear, the
			//    validation will be taken care of then.
			rowIsNested := ctr.row != 0
			rowIsAlone := ctx.st.Len() < 1
			malformedRow := rowIsNested || rowIsAlone
			if malformedRow {
				log.Println(errMalformedRowMsg)
				return errors.New(errMalformedRowMsg)
			}

			if ctr.row == 0 {
				log.Println("Adding row to table")
				*tableObj = append(*tableObj, *rowObj)
				*rowObj = *new(row)
			}
		case atom.Th:
			log.Println("Encountered Header Column")
			fallthrough
		case atom.Td:
			log.Println("Reducing Column")
			ctr.col--
			// There should never be the following cases:
			//  * Nested Columns: <td><td></td></td>
			//  * Standalone Columns. The column must be preceded by
			//    at least two elements, which is expected to be
			//    <table> and <tr>
			//
			//    No validation is performed on the last case at this
			//    stage. When those token's closing tags appear, the
			//    validation will be taken care of then.
			colIsNested := ctr.col != 0
			colIsInWrongPos := ctx.st.Len() < 2
			malformedCol := colIsNested || colIsInWrongPos
			if malformedCol {
				log.Println(errMalformedColMsg)
				return errors.New(errMalformedColMsg)
			}

			log.Println("Adding column data to row")
			*rowObj = append(*rowObj, ctx.sb.String())
			ctx.sb.Reset()
		default:
			err := fmt.Errorf(errUnexpectedTokenMsgFmt, tK.DataAtom)
			log.Println(err)
			return err
		}
		return nil
	} else {
		log.Println("Invalid Closing Tag")
		return fmt.Errorf(errMismatchedTokenMsgFmt, tK.DataAtom, token.DataAtom)
	}

	panic(errImpossibleCaseMsg)
}
