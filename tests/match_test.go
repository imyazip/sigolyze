package tests

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/imyazip/sigolyze"
)

var (
	compiledSignatures *sigolyze.Compiler
	once               sync.Once
)

func generateSignatures(count int) [][]byte {
	rand.New(rand.NewSource(42))

	var signatures [][]byte
	for i := 0; i < count; i++ {
		// Случайное количество паттернов от 1 до 10
		numPatterns := rand.Intn(10) + 1
		patterns := `{"name": "test_signature", "patterns": [`
		for j := 0; j < numPatterns; j++ {
			patterns += fmt.Sprintf(`{"name": "pattern%d_%d", "value": "test%d_%d", "is_regex": false}`, i, j, i, j)
			if j < numPatterns-1 {
				patterns += ","
			}
		}
		patterns += `], `

		// Случайное количество тегов от 1 до 5
		numTags := rand.Intn(5) + 1
		tags := `"tags": [`
		for k := 0; k < numTags; k++ {
			tags += fmt.Sprintf(`"tag%d_%d"`, i, k)
			if k < numTags-1 {
				tags += ","
			}
		}
		tags += `], `

		patterns += tags + `"meta": []}`

		signatures = append(signatures, []byte(patterns))
	}
	return signatures
}

func initCompiler() *sigolyze.Compiler {
	once.Do(func() {
		compiler := sigolyze.NewCompiler()
		signs := generateSignatures(10000)
		for _, sign := range signs {
			compiler.LoadSignature(sign)
		}
		compiledSignatures = compiler
	})
	return compiledSignatures
}

func TestMatch(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.Match(compiler, "Value1")

	if matches[0] != &compiler.Signatures[0] {
		t.Errorf("Failed matching")
	}
}

func TestMatchAho(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.MatchAho(compiler, "Value1")

	if matches[0] != &compiler.Signatures[0] {
		t.Errorf("Failed matching")
	}
}

func TestMatchTags(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.MatchTags(compiler, "Value1 Value2 Value3 Value12", []string{"tag1"})

	if matches[0].Name != compiler.Signatures[0].Name {
		t.Errorf("Failed matching")
	}
}

func TestMatchTagsAho(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.MatchTagsAho(compiler, "Value1", []string{"tag1"})

	if matches[0].Name != compiler.Signatures[0].Name {
		t.Errorf("Failed matching")
	}
}

func BenchmarkMatch(b *testing.B) {
	compiler := initCompiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sigolyze.Match(compiler, "Value1 Value2 Value3 Value12")
	}
}

func BenchmarkMatchAho(b *testing.B) {
	compiler := initCompiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sigolyze.MatchAho(compiler, "Value1 Value2 Value3 Value12")
	}
}

func BenchmarkMatchTags(b *testing.B) {
	compiler := initCompiler()
	data := "Value1 Value2 Value3 Value12"
	tags := []string{"tag1, tag2"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sigolyze.MatchTags(compiler, data, tags)
	}
}

func BenchmarkMatchTagsAho(b *testing.B) {
	compiler := initCompiler()
	data := "Value1 Value2 Value3 Value12"
	tags := []string{"tag1, tag2"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sigolyze.MatchTagsAho(compiler, data, tags)
	}
}
