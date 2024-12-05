package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/imyazip/sigolyze"
)

func generateSignatures(count int) [][]byte {
	rand.Seed(time.Now().UnixNano()) // Инициализация случайного генератора

	var signatures [][]byte
	for i := 0; i < count; i++ {
		numPatterns := rand.Intn(10) + 1 // Случайное количество паттернов от 1 до 10
		patterns := `{"name": "test_signature", "patterns": [`
		for j := 0; j < numPatterns; j++ {
			patterns += fmt.Sprintf(`{"name": "pattern%d_%d", "value": "test%d_%d", "is_regex": false}`, i, j, i, j)
			if j < numPatterns-1 {
				patterns += ","
			}
		}
		patterns += `], "tags": ["tag1"], "meta": []}`
		signatures = append(signatures, []byte(patterns))
	}
	return signatures
}

func setupCompiler(numSignatures int) *sigolyze.Compiler {
	compiler := sigolyze.NewCompiler()
	signatures := generateSignatures(numSignatures)

	for _, sig := range signatures {
		_ = compiler.LoadSignature(sig) // Игнорируем ошибки для простоты
	}

	return compiler
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

	matches := sigolyze.MatchTags(compiler, "Value1", []string{"tag1"})

	if matches[0] != &compiler.Signatures[0] {
		t.Errorf("Failed matching")
	}
}

func TestMatchTagsAho(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignatureFromJson("example.json")

	matches := sigolyze.MatchTagsAho(compiler, "Value1", []string{"tag1"})

	if matches[0] != &compiler.Signatures[0] {
		t.Errorf("Failed matching")
	}
}

func BenchmarkMatch(b *testing.B) {
	compiler := sigolyze.NewCompiler()
	signs := generateSignatures(100)
	for _, sign := range signs {
		compiler.LoadSignature(sign)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sigolyze.Match(compiler, "Value1")
	}
	b.StopTimer()
}

func BenchmarkMatchAho(b *testing.B) {
	compiler := sigolyze.NewCompiler()
	signs := generateSignatures(100)
	for _, sign := range signs {
		compiler.LoadSignature(sign)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sigolyze.MatchAho(compiler, "Value1")
	}
	b.StopTimer()
}

func BenchmarkMatchTags(b *testing.B) {
	compiler := sigolyze.NewCompiler()
	signs := generateSignatures(100)
	for _, sign := range signs {
		compiler.LoadSignature(sign)
	}

	data := "test0_1 test10_3 test20_2" // Данные для поиска
	tags := []string{"tag1"}            // Теги для фильтрации

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sigolyze.MatchTags(compiler, data, tags)
	}
}

func BenchmarkMatchTagsAho(b *testing.B) {
	compiler := sigolyze.NewCompiler()
	signs := generateSignatures(100)
	for _, sign := range signs {
		compiler.LoadSignature(sign)
	}

	data := "test0_1 test10_3 test20_2" // Данные для поиска
	tags := []string{"tag1"}            // Теги для фильтрации

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sigolyze.MatchTagsAho(compiler, data, tags)
	}
}
