package birthday_test

import (
	"strings"
	"testing"

	"github.com/keep94/birthday"
	"github.com/keep94/toolbox/date_util"
	asserts "github.com/stretchr/testify/assert"
)

func TestMalformattedLine(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat
`
	search := birthday.NewSearch("")
	err := birthday.Read(strings.NewReader(fileContents), search)
	assert.EqualError(err, "Line 4 malformatted")
}

func TestBadDateWithYear(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
	# This is a comment

Jack Sprat	08/32/2006
`
	search := birthday.NewSearch("")
	err := birthday.Read(strings.NewReader(fileContents), search)
	assert.EqualError(err, "Line 4 contains invalid birthday")
}

func TestBadDate(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat	08/32
`
	search := birthday.NewSearch("")
	err := birthday.Read(strings.NewReader(fileContents), search)
	assert.EqualError(err, "Line 4 contains invalid birthday")
}

func TestReadLines(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

Jack Sprat	08/31/2006	Tea
Alice Doe	12/15
`
	search := birthday.NewSearch("")
	err := birthday.Read(strings.NewReader(fileContents), search)
	assert.NoError(err)
	assert.Equal([]birthday.Entry{
		{
			Name:     "Alice Doe",
			Birthday: date_util.YMD(0, 12, 15),
		},
		{
			Name:     "Jack Sprat",
			Birthday: date_util.YMD(2006, 8, 31),
		},
	}, search.Results())
}

func TestReadLinesWithWhitespace(t *testing.T) {
	assert := asserts.New(t)

	fileContents := `
# This is a comment

	Jack Sprat	08/31/2006	Tea
	Alice Doe	12/15     
`
	search := birthday.NewSearch("")
	err := birthday.Read(strings.NewReader(fileContents), search)
	assert.NoError(err)
	assert.Equal([]birthday.Entry{
		{
			Name:     "Alice Doe",
			Birthday: date_util.YMD(0, 12, 15),
		},
		{
			Name:     "Jack Sprat",
			Birthday: date_util.YMD(2006, 8, 31),
		},
	}, search.Results())
}
