package tests

import (
	"testing"

	"github.com/imyazip/sigolyze"
)

var signs string = `
{
    "name": "Example Signature",
    "patterns": [
      {
        "name": "Pattern 1",
        "value": "Value1",
        "is_regex": true
      },
      {
        "name": "Pattern 2",
        "value": "Value2",
        "is_regex": false
      }
    ],
    "tags": ["tag1", "tag2"],
    "meta": [
      {
        "name": "Meta1",
        "info": ["detail1", "detail2"]
      },
      {
        "name": "Meta2",
        "info": ["detail3"]
      }
    ]
  }
	`

func TestLoadSignature(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadSignature([]byte(signs))
	expectedPatterns := []sigolyze.Pattern{
		{
			Name:    "Pattern 1",
			Value:   "Value1",
			IsRegex: true,
		},
		{
			Name:    "Pattern 2",
			Value:   "Value2",
			IsRegex: false,
		},
	}
	expectedMeta := []sigolyze.MetaInfo{
		{Name: "Meta1",
			Info: []string{"detail1", "detail2"},
		},
		{Name: "Meta2",
			Info: []string{"detail3"},
		},
	}

	for index := range expectedPatterns {
		if expectedPatterns[index].Name != compiler.Signatures[0].Patterns[index].Name {
			t.Errorf("Failed loading pattern names")
		}

		if expectedPatterns[index].Value != compiler.Signatures[0].Patterns[index].Value {
			t.Errorf("Failed loading pattern values")
		}

		if expectedPatterns[index].IsRegex != compiler.Signatures[0].Patterns[index].IsRegex {
			t.Errorf("Failed loading pattern values")
		}
	}

	for index := range expectedMeta {
		if expectedMeta[index].Name != compiler.Signatures[0].Meta[index].Name {
			t.Error("Failed loading metadata names")
		}

		for metaIndex := range expectedMeta[index].Info {
			if expectedMeta[index].Info[metaIndex] != compiler.Signatures[0].Meta[index].Info[metaIndex] {
				t.Error("Failed loading metadata info")
			}
		}
	}
}

func TestLoadSignatureMultiple(t *testing.T) {
	compiler := sigolyze.NewCompiler()

	// Несколько тестовых сигнатур
	signaturesJSON := [][]byte{
		[]byte(`{
			"name": "signature1",
			"patterns": [
				{"name": "pattern1", "value": "test1", "is_regex": false}
			],
			"tags": ["tag1"],
			"meta": [
				{"name": "meta1", "info": ["info1"]}
			]
		}`),
		[]byte(`{
			"name": "signature2",
			"patterns": [
				{"name": "pattern2", "value": "test2", "is_regex": false},
				{"name": "pattern3", "value": "test3", "is_regex": false}
			],
			"tags": [],
			"meta": []
		}`),
		[]byte(`{
			"name": "signature3",
			"patterns": [],
			"tags": ["tag2"],
			"meta": [
				{"name": "meta2", "info": ["info2", "info3"]}
			]
		}`),
	}

	// Загружаем сигнатуры
	for _, sigJSON := range signaturesJSON {
		err := compiler.LoadSignature(sigJSON)
		if err != nil {
			t.Fatalf("Failed to load signature: %v", err)
		}
	}

	// Проверяем количество загруженных сигнатур
	if len(compiler.Signatures) != len(signaturesJSON) {
		t.Fatalf("Expected %d signatures, got %d", len(signaturesJSON), len(compiler.Signatures))
	}

	// Проверяем содержимое первой сигнатуры
	sign1 := compiler.Signatures[0]
	if sign1.Name != "signature1" || len(sign1.Patterns) != 1 || sign1.Patterns[0].Value != "test1" {
		t.Errorf("Signature 1 mismatch: %+v", sign1)
	}

	// Проверяем содержимое второй сигнатуры
	sign2 := compiler.Signatures[1]
	if sign2.Name != "signature2" || len(sign2.Patterns) != 2 || sign2.Patterns[1].Value != "test3" {
		t.Errorf("Signature 2 mismatch: %+v", sign2)
	}

	// Проверяем содержимое третьей сигнатуры
	sign3 := compiler.Signatures[2]
	if sign3.Name != "signature3" || len(sign3.Patterns) != 0 || sign3.Meta[0].Name != "meta2" {
		t.Errorf("Signature 3 mismatch: %+v", sign3)
	}

	// Проверяем, что матчер корректно ищет строки
	matches := compiler.Signatures[0].Matcher.Match([]byte("this is a test1 string"))
	if len(matches) == 0 {
		t.Error("Matcher failed to find 'test1' in signature 1")
	}

	matches = compiler.Signatures[1].Matcher.Match([]byte("this is a test3 string"))
	if len(matches) == 0 {
		t.Error("Matcher failed to find 'test3' in signature 2")
	}
}
