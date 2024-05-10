package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name string
		input string
		want  Status
	} {
		{"Pending status", "pending", PENDING},
		{"In progress status", "in_progress", IN_PROGRESS},
		{"Invalid existing status", "In Progress", ""},
		{"Non-existent status", "Bogus", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := ParseStatus(test.input)
			assert.Equal(test.want, actual, "Got %s, want %s", actual, test.want)
		})
	}
}

