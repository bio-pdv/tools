package parse

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

const (
	blankCol           = "blank col"
	testSuffix         = "test_data_"
	testSuffix0        = testSuffix + "0"
	testSuffix1        = testSuffix + "1"
	testEvidence       = "AB"
	testSeqId          = "AB0123456"
	testPosition       = "12,345"
	testMutation       = "+G"
	testFrequency      = "12%"
	testAnnotation     = "abcdefg&nbsp;(&#1234;567/+89)"
	testGene           = "<i>ABC0123</i>&nbsp;&larr;&nbsp;/&nbsp;&larr;&nbsp;<i>ABC5678</i>"
	expectedTestGene   = "ABC0123&nbsp;&larr;&nbsp;/&nbsp;&larr;&nbsp;ABC5678"
	testDescription    = "abcedfghi/jklmnopqrst"
	validBreseq027Html = `<!DOCTYPE html
PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" lang="en" xml:lang="en">
<html>
<head>
<title>BRESEQ :: Mutation Predictions</title>
<style type = "text/css">
body {font-family: sans-serif; font-size: 11pt;}
th {background-color: rgb(0,0,0); color: rgb(255,255,255);}
table {background-color: rgb(1,0,0); color: rgb(0,0,0);}
tr {background-color: rgb(255,255,255);}
.mutation_in_codon {color:red; text-decoration : underline;}

</style>
<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1" />
<script type="text/javascript">
  function hideTog(divID) {
    var item = document.getElementById(divID);
    if (item) {
      item.className=(item.className=='hidden')?'unhidden':'hidden';
    }
  }
</script>

</head>
<body>
<table width="100%" border="0" cellspacing="0" cellpadding="3">
<tr>
<td><a href="http://barricklab.org/breseq"><img src="evidence/breseq_small.png" /></a></td>
<td width="100%">
<b><i>breseq</i></b>&nbsp;&nbsp;version 0.27.1&nbsp;&nbsp;revision 87c22d663cc3
<br><a href="index.html">mutation predictions</a> | 
<a href="marginal.html">marginal predictions</a> | 
<a href="summary.html">summary statistics</a> | 
<a href="output.gd">genome diff</a> | 
<a href="log.txt">command line log</a>
</td></tr></table>

<p>
<!--Mutation Predictions -->
<p>
<!--Output Html_Mutation_Table_String-->
<table border="0" cellspacing="1" cellpadding="3">
<tr><th colspan="8" align="left" class="mutation_header_row">Predicted mutations</th></tr><tr>
<th>evidence</th>
<th>seq&nbsp;id</th>
<th>position</th>
<th>mutation</th>
<th>freq</th>
<th>annotation</th>
<th>gene</th>
<th width="100%">description</th>
</tr>

<!-- Item Lines -->

<!-- Print The Table Row -->
<tr class="normal_table_row">
<td align="center"><a href="evidence/INS_0.html">` + testEvidence + testSuffix0 + `</a></td><!-- Evidence -->
<td align="center">` + testSeqId + testSuffix0 + `</td><!-- Seq_Id -->
<td align="right">` + testPosition + testSuffix0 + `</td><!-- Position -->
<td align="center">` + testMutation + testSuffix0 + `</td><!-- Cell Mutation -->
<td align="right">` + testFrequency + testSuffix0 + `</td>
<td align="center">` + testAnnotation + testSuffix0 + `</td>
<td align="center">` + testGene + testSuffix0 + `</td>
<td align="left">` + testDescription + testSuffix0 + `</td>
</tr>
<!-- End Table Row -->

<!-- Print The Table Row -->
<tr class="normal_table_row">
<td align="center"><a href="evidence/INS_0.html">` + testEvidence + testSuffix1 + `</a></td><!-- Evidence -->
<td align="center">` + testSeqId + testSuffix1 + `</td><!-- Seq_Id -->
<td align="right">` + testPosition + testSuffix1 + `</td><!-- Position -->
<td align="center">` + testMutation + testSuffix1 + `</td><!-- Cell Mutation -->
<td align="right">` + testFrequency + testSuffix1 + `</td>
<td align="center">` + testAnnotation + testSuffix1 + `</td>
<td align="center">` + testGene + testSuffix1 + `</td>
<td align="left">` + testDescription + testSuffix1 + `</td>
</tr>
<!-- End Table Row -->
</table>
</body>
</html>`
)

var (
	dataRow = row{
		"evid",
		"AB_012345",
		"12,456",
		"+G",
		"100%",
		"abcdefg&nbsp;(&#1234;567/+89)",
		"ABC0123</i>&nbsp;&larr;&nbsp;/&nbsp;&larr;&nbsp;DV4567",
		"abcdefghijklmnopqrstuvwxyz",
	}
	throwAwayRow = row{"throw away row"}
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestParseBreseq027HtmlFile(t *testing.T) {
	testReader := strings.NewReader(validBreseq027Html)
	testResults, testErr := parseBreseq027HtmlFile(testReader)
	assert.Nil(t, testErr)
	assert.NotNil(t, testResults)
	assert.Equal(t, 1, len(testResults))
	assert.Equal(t, 2, len(testResults[0]))

	for i, rowRes := range testResults[0] {
		iStr := fmt.Sprintf("%d", i)
		assert.Equal(t, html.UnescapeString(testSeqId+testSuffix+iStr), rowRes.SequenceId)
		assert.Equal(t, html.UnescapeString(testPosition+testSuffix+iStr), rowRes.Position)
		assert.Equal(t, html.UnescapeString(testMutation+testSuffix+iStr), rowRes.Mutation)
		assert.Equal(t, html.UnescapeString(testFrequency+testSuffix+iStr), rowRes.Frequency)
		assert.Equal(t, html.UnescapeString(testAnnotation+testSuffix+iStr), rowRes.Annotation)
		assert.Equal(t, html.UnescapeString(expectedTestGene+testSuffix+iStr), rowRes.Gene)
		assert.Equal(t, html.UnescapeString(testDescription+testSuffix+iStr), rowRes.Description)
	}
}

func TestChangeBreseq027TableToSeqAnnotation(t *testing.T) {
	testTable := table{
		row{"throw away header"},
		row(breseq027HtmlDataHeaders),
		dataRow,
	}
	testSaTable := changeBreseq027TableToSeqAnnotation(testTable)
	assert.Equal(t, 1, len(testSaTable))

	testSa := testSaTable[0]
	assert.Equal(t, dataRow[1], testSa.SequenceId)
	assert.Equal(t, dataRow[2], testSa.Position)
	assert.Equal(t, dataRow[3], testSa.Mutation)
	assert.Equal(t, dataRow[4], testSa.Frequency)
	assert.Equal(t, dataRow[5], testSa.Annotation)
	assert.Equal(t, dataRow[6], testSa.Gene)
	assert.Equal(t, dataRow[7], testSa.Description)
	assert.Equal(t, testSa.Application, string(breseq))
	assert.Equal(t, testSa.AppVersion, string(breseqVers027Number))
	assert.Equal(t, testSa.Generation, "")
}

func TestChangeBreseq027TableToSeqAnnotationNoHeader(t *testing.T) {
	testTable := table{
		dataRow,
	}

	testSaTable := changeBreseq027TableToSeqAnnotation(testTable)
	assert.Nil(t, testSaTable)
}

func TestIsBreseq027VersTable(t *testing.T) {
	cases := []struct {
		name      string
		testTable table
	}{
		{
			name: "Appended Patch",
			testTable: table{
				row{blankCol, breseqVers027Prefix + ".1"},
			},
		},
		{
			name: "Non-breaking Space",
			testTable: table{
				row{blankCol, "breseq\u00A0\u00A0version\u00A0\u00A00.27"},
			},
		},
		{
			name: "New-line Character",
			testTable: table{
				row{blankCol, "breseq\nversion\n0.27"},
			},
		},
		{
			name: "Spaces",
			testTable: table{
				row{blankCol, "breseq version 0.27 .1"},
			},
		},
		{
			name: "All Combined Cases",
			testTable: table{
				row{blankCol, "breseq\u00A0\n\u00A0 version\u00A0\u00A0\n 0.27.1"},
			},
		},
	}

	for _, c := range cases {
		assert.True(t, isBreseq027VersTable(c.testTable), c.name)
	}
}

func TestIsBreseq027InvalidVersTable(t *testing.T) {
	cases := []struct {
		name      string
		testTable table
	}{
		{
			name:      "Empty",
			testTable: table{},
		},
		{
			name: "Different",
			testTable: table{
				row{blankCol, "breseq version 0.28"},
			},
		},
		{
			name: "Invalid",
			testTable: table{
				row{blankCol, "invalid version"},
			},
		},
		{
			name: "No Blank Column",
			testTable: table{
				row{"breseq version 0.27"},
			},
		},
		{
			name: "Blank Column Only",
			testTable: table{
				row{blankCol},
			},
		},
	}

	for _, c := range cases {
		assert.False(t, isBreseq027VersTable(c.testTable), c.name)
	}
}

func TestIsBreseq027DataTable(t *testing.T) {
	dataTable := table{
		row{"throw away row"},
		row(breseq027HtmlDataHeaders),
		row{"test row"},
	}

	assert.True(t, isBreseq027DataTable(dataTable))
}

func TestIsBreseq027InvalidDataTable(t *testing.T) {
	cases := []struct {
		name      string
		testTable table
	}{
		{
			name:      "Empty",
			testTable: table{},
		},
		{
			name: "Invalid Header",
			testTable: table{
				throwAwayRow,
				row{"invalid header"},
				row{"data"},
			},
		},
		{
			name: "Less Than 2 Rows",
			testTable: table{
				row{"random row"},
			},
		},
		{
			name: "Mismatched Headers",
			testTable: table{
				throwAwayRow,
				row(append(breseq027HtmlDataHeaders[3:], breseq027HtmlDataHeaders[:3]...)),
				row{"data"},
			},
		},
		{
			name: "Not Enough Headers",
			testTable: table{
				throwAwayRow,
				row(breseq027HtmlDataHeaders[:4]),
				row{"data"},
			},
		},
		{
			name: "Empty Headers",
			testTable: table{
				throwAwayRow,
				row{},
				row{"data"},
			},
		},
	}

	for _, c := range cases {
		assert.False(t, isBreseq027DataTable(c.testTable), c.name)
	}
}

func TestIsBreseq027(t *testing.T) {
	versTable := table{
		row{blankCol, breseqVers027Prefix + ".12345"},
	}
	headerTable := table{
		throwAwayRow,
		row(breseq027HtmlDataHeaders),
	}
	testTables := []table{
		versTable,
		headerTable,
	}

	assert.True(t, isBreseq027(testTables))
}

func TestIsBreseq027NotEnoughTables(t *testing.T) {
	notEnoughTables := []table{
		table{
			row{"blank col", breseqVers027Prefix + ".12345"},
		},
	}

	versTable := table{
		row{blankCol, breseqVers027Prefix + ".12345"},
	}
	headerTable := table{
		throwAwayRow,
		row(breseq027HtmlDataHeaders),
	}
	swappedTables := []table{
		headerTable,
		versTable,
	}

	assert.False(t, isBreseq027(notEnoughTables))
	assert.False(t, isBreseq027(swappedTables))
}
