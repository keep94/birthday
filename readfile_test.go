package birthday_test

import (
	"strings"
	"testing"

	"github.com/keep94/birthday"
	"github.com/keep94/consume2"
	"github.com/keep94/toolbox/date_util"
	asserts "github.com/stretchr/testify/assert"
)

func TestMalformattedLine(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents), consume2.AppendTo(&entries))
	assert.EqualError(err, "Line 4 malformatted")
}

func TestBadDateWithYear(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
	# This is a comment

Jack Sprat	08/32/2006
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents), consume2.AppendTo(&entries))
	assert.EqualError(err, "Line 4 contains invalid birthday")
}

func TestBadDate(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat	08/32
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents), consume2.AppendTo(&entries))
	assert.EqualError(err, "Line 4 contains invalid birthday")
}

func TestReadLines(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat	08/31/2006	Tea
Alice Doe	12/15
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents), consume2.AppendTo(&entries))
	assert.NoError(err)
	assert.Equal([]birthday.Entry{
		{
			Name:     "Jack Sprat",
			Birthday: date_util.YMD(2006, 8, 31),
		},
		{
			Name:     "Alice Doe",
			Birthday: date_util.YMD(0, 12, 15),
		},
	}, entries)
}

func TestReadLinesQuitEarly(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat	08/31/2006	Tea
Alice Doe	12/15
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents),
		consume2.Slice(consume2.AppendTo(&entries), 0, 1))
	assert.NoError(err)
	assert.Equal([]birthday.Entry{
		{
			Name:     "Jack Sprat",
			Birthday: date_util.YMD(2006, 8, 31),
		},
	}, entries)
}

func TestReadLinesWithWhitespace(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

	Jack Sprat	08/31/2006	Tea
	Alice Doe	12/15     
`
	var entries []birthday.Entry
	err := birthday.Read(
		strings.NewReader(fileContents), consume2.AppendTo(&entries))
	assert.NoError(err)
	assert.Equal([]birthday.Entry{
		{
			Name:     "Jack Sprat",
			Birthday: date_util.YMD(2006, 8, 31),
		},
		{
			Name:     "Alice Doe",
			Birthday: date_util.YMD(0, 12, 15),
		},
	}, entries)
}
