package tests

import (
	"testing"

	"github.com/imyazip/sigolyze"
)

func TestMatch(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.Match(compiler, "Value1")

	if matches[0] != &compiler.Signatures[0] {
		t.Errorf("Failed matching")
	}
}
