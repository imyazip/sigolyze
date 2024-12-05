package tests

import (
	"testing"

	"github.com/imyazip/sigolyze"
)

func TestLoadRulesFromJson(t *testing.T) {
	compiler := sigolyze.NewCompiler()
	compiler.LoadRulesFromJson("example.json")
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
		if expectedPatterns[index].Name != compiler.Signatures.Patterns[index].Name {
			t.Errorf("Failed loading pattern names")
		}

		if expectedPatterns[index].Value != compiler.Signatures.Patterns[index].Value {
			t.Errorf("Failed loading pattern values")
		}

		if expectedPatterns[index].IsRegex != compiler.Signatures.Patterns[index].IsRegex {
			t.Errorf("Failed loading pattern values")
		}
	}

	for index := range expectedMeta {
		if expectedMeta[index].Name != compiler.Signatures.Meta[index].Name {
			t.Error("Failed loading metadata names")
		}

		for metaIndex := range expectedMeta[index].Info {
			if expectedMeta[index].Info[metaIndex] != compiler.Signatures.Meta[index].Info[metaIndex] {
				t.Error("Failed loading metadata info")
			}
		}
	}
}
