package parse

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"testing"
)

const (
	testVal                   = "value"
	itTestVal                 = "italicized value"
	wellFormedHtmlTableString = `<table><tr><td>` + testVal + ` <i>` + itTestVal + `</i></td></tr></table>`
)

type mockStack struct {
	mock.Mock
}

func (m *mockStack) Push(arg interface{}) {
	m.Called(arg)
}

func (m *mockStack) Pop() interface{} {
	return m.Called()[0]
}

func (m *mockStack) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *mockStack) Peek() interface{} {
	return m.Called()[0]
}

type mockTokenizer struct {
	mock.Mock
}

func (m *mockTokenizer) Next() html.TokenType {
	args := m.Called()
	return args[0].(html.TokenType)
}

func (m *mockTokenizer) Token() html.Token {
	args := m.Called()
	return args[0].(html.Token)
}

func (m *mockTokenizer) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTokenizer) Raw() []byte {
	args := m.Called()
	return args[0].([]byte)
}

func TestParseDataTableHtmlTokenizer(t *testing.T) {
	testReader := strings.NewReader(wellFormedHtmlTableString)
	testResults, testErr := parseDataTableHtmlTokenizer(testReader)
	assert.Nil(t, testErr)
	assert.NotNil(t, testResults)
	assert.Equal(t, 1, len(testResults))
	testTable := testResults[0]
	assert.Equal(t, 1, len(testTable))
	testRow := testTable[0]
	assert.Equal(t, 1, len(testRow))
	testColVal := testRow[0]
	assert.Equal(t, testVal+" "+itTestVal, testColVal)
}

func TestHandleTextTagToken(t *testing.T) {
	testString := "Test String"
	mockSt, mTokenizer := new(mockStack), new(mockTokenizer)
	colStartToken := &html.Token{
		Type:     html.StartTagToken,
		DataAtom: atom.Td,
	}
	testToken := html.Token{
		Type: html.TextToken,
		Data: testString,
	}
	mockSt.On("Peek").Return(colStartToken).Once()
	mockSt.On("Len").Return(1).Once()
	mTokenizer.On("Token").Return(testToken).Once()
	testCtx := &parserContext{
		sb: strings.Builder{},
		st: mockSt,
		tk: mTokenizer,
	}
	testErr := handleTextTagToken(testCtx)
	assert.Nil(t, testErr)
	assert.Equal(t, testString, testCtx.sb.String())
	mockSt.AssertCalled(t, "Peek")
	mTokenizer.AssertCalled(t, "Token")
}

func TestHandleTextTagTokenSkip(t *testing.T) {
	mockSt, mTokenizer := new(mockStack), new(mockTokenizer)
	mockSt.On("Len").Return(0).Once()
	testCtx := &parserContext{
		st: mockSt,
		tk: mTokenizer,
	}

	assert.Nil(t, handleTextTagToken(testCtx))
	mockSt.AssertCalled(t, "Len")
}

func TestHandleStartTagToken(t *testing.T) {
	ctr := &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	eCtr := &tagCtr{
		col:   0,
		row:   0,
		table: 1,
	}
	assertHandleStartTagToken(t, atom.Table, true, ctr, eCtr)

	eCtr = &tagCtr{
		col:   0,
		row:   1,
		table: 1,
	}
	assertHandleStartTagToken(t, atom.Tr, true, ctr, eCtr)

	eCtr = &tagCtr{
		col:   1,
		row:   1,
		table: 1,
	}
	assertHandleStartTagToken(t, atom.Td, true, ctr, eCtr)

	eCtr = &tagCtr{
		col:   2,
		row:   1,
		table: 1,
	}
	assertHandleStartTagToken(t, atom.Th, true, ctr, eCtr)
}

func TestHandleStartTagTokenSkipTokens(t *testing.T) {
	ctr := &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	eCtr := &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	assertHandleStartTagToken(t, atom.A, false, ctr, eCtr)
	assertHandleStartTagToken(t, atom.Br, false, ctr, eCtr)
}

func assertHandleStartTagToken(t *testing.T, tag atom.Atom, stackPushExpected bool, ctr *tagCtr, eCtr *tagCtr) {
	testToken := html.Token{
		DataAtom: tag,
	}
	mockSt, mTokenizer := new(mockStack), new(mockTokenizer)
	mockSt.On("Push", mock.Anything)
	mTokenizer.On("Token").Return(testToken).Once()

	ctx := &parserContext{
		sb: strings.Builder{},
		st: mockSt,
		tk: mTokenizer,
	}
	handleStartTagToken(ctr, ctx)
	assert.Equal(t, eCtr.col, ctr.col)
	assert.Equal(t, eCtr.row, ctr.row)
	assert.Equal(t, eCtr.table, ctr.table)
	mTokenizer.AssertCalled(t, "Token")
	mTokenizer.AssertNumberOfCalls(t, "Token", 1)
	mockSt.AssertNotCalled(t, "Pop")
	if stackPushExpected {
		mockSt.AssertCalled(t, "Push", mock.Anything)
		mockSt.AssertNumberOfCalls(t, "Push", 1)
	} else {
		mockSt.AssertNotCalled(t, "Push", mock.Anything)
	}
}

func TestHandleEndTagToken(t *testing.T) {
	testString := "Test String"
	ctr := &tagCtr{
		col:   1,
		row:   1,
		table: 1,
	}
	tRow, tTable := new(row), &table{row{"sentinel"}}
	tResults := new([]table)
	mockSt, mTokenizer := new(mockStack), new(mockTokenizer)
	ctx := &parserContext{
		sb: strings.Builder{},
		st: mockSt,
		tk: mTokenizer,
	}
	testToken := &html.Token{
		DataAtom: atom.Td,
	}

	ctx.sb.WriteString(testString)
	mTokenizer.On("Token").Return(*testToken).Once()
	mockSt.On("Pop").Return(testToken).Once()
	mockSt.On("Len").Return(2).Once()
	testErr := handleEndTagToken(tRow, tTable, tResults, ctr, ctx)
	mTokenizer.AssertCalled(t, "Token")
	mockSt.AssertCalled(t, "Pop")
	assert.Nil(t, testErr)
	assert.Equal(t, 0, ctr.col)
	assert.Equal(t, 1, ctr.row)
	assert.Equal(t, 1, ctr.table)
	assert.Equal(t, 0, len(*tResults))
	assert.Equal(t, 1, len(*tTable))

	mockSt, mTokenizer = new(mockStack), new(mockTokenizer)
	ctx.st = mockSt
	ctx.tk = mTokenizer
	testToken.DataAtom = atom.Tr
	mTokenizer.On("Token").Return(*testToken).Once()
	mockSt.On("Pop").Return(testToken).Once()
	mockSt.On("Len").Return(1).Once()
	testErr = handleEndTagToken(tRow, tTable, tResults, ctr, ctx)
	mTokenizer.AssertCalled(t, "Token")
	mockSt.AssertCalled(t, "Pop")
	assert.Nil(t, testErr)
	assert.Equal(t, 0, ctr.col)
	assert.Equal(t, 0, ctr.row)
	assert.Equal(t, 1, ctr.table)
	assert.Equal(t, 0, len(*tResults))
	assert.Equal(t, 2, len(*tTable))

	mockSt, mTokenizer = new(mockStack), new(mockTokenizer)
	ctx.st = mockSt
	ctx.tk = mTokenizer
	testToken.DataAtom = atom.Table
	mTokenizer.On("Token").Return(*testToken).Once()
	mockSt.On("Pop").Return(testToken).Once()
	mockSt.On("Len").Return(0).Twice()
	testErr = handleEndTagToken(tRow, tTable, tResults, ctr, ctx)
	mTokenizer.AssertCalled(t, "Token")
	mockSt.AssertCalled(t, "Pop")
	assert.Nil(t, testErr)
	assert.Equal(t, 0, ctr.col)
	assert.Equal(t, 0, ctr.row)
	assert.Equal(t, 0, ctr.table)
	assert.Equal(t, 1, len(*tResults))
	assert.Equal(t, 0, len(*tTable))

	assert.NotEqual(t, (*tResults)[0], *tTable)
	assert.Equal(t, "", ctx.sb.String())

	foundTestString := false
	tTableRes := (*tResults)[0]
	for i := 0; i < len(tTableRes) && !foundTestString; i++ {
		for _, val := range tTableRes[i] {
			if val == testString {
				foundTestString = true
				break
			}
		}
	}
	assert.True(t, foundTestString)
}

func TestHandleEndTagTokenMalformedTable(t *testing.T) {
	ctr := &tagCtr{
		col:   0,
		row:   0,
		table: 1,
	}
	eCtr := &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Table, 1, ctr, eCtr, errMalformedTableMsg)
	ctr = &tagCtr{
		col:   0,
		row:   0,
		table: 2,
	}
	eCtr = &tagCtr{
		col:   0,
		row:   0,
		table: 1,
	}
	assertHandleEndTagTokenMalformed(t, atom.Table, 0, ctr, eCtr, errMalformedTableMsg)
	ctr = &tagCtr{
		col:   0,
		row:   0,
		table: -1,
	}
	eCtr = &tagCtr{
		col:   0,
		row:   0,
		table: -2,
	}
	assertHandleEndTagTokenMalformed(t, atom.Table, 0, ctr, eCtr, errMalformedTableMsg)
}

func TestHandleEndTagTokenMalformedRow(t *testing.T) {
	ctr := &tagCtr{
		col:   0,
		row:   -1,
		table: 0,
	}
	eCtr := &tagCtr{
		col:   0,
		row:   -2,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Tr, 0, ctr, eCtr, errMalformedRowMsg)
	ctr = &tagCtr{
		col:   0,
		row:   2,
		table: 0,
	}
	eCtr = &tagCtr{
		col:   0,
		row:   1,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Tr, 0, ctr, eCtr, errMalformedRowMsg)
	ctr = &tagCtr{
		col:   0,
		row:   1,
		table: 0,
	}
	eCtr = &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Tr, 0, ctr, eCtr, errMalformedRowMsg)
}

func TestHandleEndTagTokenMalformedCol(t *testing.T) {
	ctr := &tagCtr{
		col:   -1,
		row:   0,
		table: 0,
	}
	eCtr := &tagCtr{
		col:   -2,
		row:   0,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Td, 0, ctr, eCtr, errMalformedColMsg)
	ctr = &tagCtr{
		col:   1,
		row:   0,
		table: 0,
	}
	eCtr = &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Td, 0, ctr, eCtr, errMalformedColMsg)
	ctr = &tagCtr{
		col:   0,
		row:   0,
		table: 0,
	}
	eCtr = &tagCtr{
		col:   -1,
		row:   0,
		table: 0,
	}
	assertHandleEndTagTokenMalformed(t, atom.Th, 0, ctr, eCtr, errMalformedColMsg)
}

func assertHandleEndTagTokenMalformed(t *testing.T, a atom.Atom, stLen int, ctr *tagCtr, eCtr *tagCtr, eErrMsg string) {
	mockSt, mTokenizer := new(mockStack), new(mockTokenizer)
	rowObj, tableObj := new(row), new(table)
	testToken := &html.Token{
		DataAtom: a,
	}
	testStackToken := &html.Token{
		DataAtom: a,
	}
	ctx := &parserContext{
		sb: strings.Builder{},
		st: mockSt,
		tk: mTokenizer,
	}
	mockSt.On("Pop").Return(testStackToken)
	mockSt.On("Len").Return(stLen).Once()
	mTokenizer.On("Token").Return(*testToken)
	results := &[]table{}
	testErr := handleEndTagToken(rowObj, tableObj, results, ctr, ctx)
	assert.NotNil(t, testErr)
	assert.Equal(t, testErr.Error(), eErrMsg)
	assert.Equal(t, eCtr.col, ctr.col)
	assert.Equal(t, eCtr.row, ctr.row)
	assert.Equal(t, eCtr.table, ctr.table)
	mockSt.AssertNotCalled(t, "Push", mock.Anything)
	mTokenizer.AssertNotCalled(t, "Next")
}
