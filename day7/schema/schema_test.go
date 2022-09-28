package schema

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	schema := Parse(&Hero{})
	fmt.Println(schema)
}
