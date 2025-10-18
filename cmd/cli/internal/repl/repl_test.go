package repl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWords(t *testing.T) {
	input := `rar race 146`
	expected := []string{"rar", "race", "146"}
	result := parseString(input)
	fmt.Println("actual:", result)
	if len(result) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(result))
	}
}

func TestParseRar(t *testing.T) {
	input := `rar "Practice JV" 146`
	expected := []string{"rar", "Practice JV", "146"}
	result := parseString(input)
	fmt.Println("actual:", result)
	assert.Equal(t, expected, result)
}
