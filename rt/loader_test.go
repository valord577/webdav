package rt

import (
	"fmt"
	"testing"
)

// @author valor.

var cachedCnf cfg

func TestLoader(t *testing.T) {
	err := ReadInFile("app.jsonc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cachedCnf)
}
