package sigolyze

import (
	"io"
	"os"
	"strings"

	"github.com/cloudflare/ahocorasick"
	json "github.com/json-iterator/go"
)

type Pattern struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"is_regex"`
}

type MetaInfo struct {
	Name string   `json:"name"`
	Info []string `json:"info"`
}

type Signature struct {
	Name     string     `json:"name"`
	Patterns []Pattern  `json:"patterns"`
	Tags     []string   `json:"tags"`
	Meta     []MetaInfo `json:"meta"`
	Matcher  *ahocorasick.Matcher
}

type Compiler struct {
	Signatures []Signature
}

func (p *Pattern) NewPattern(name string, value string, isRegex bool) *Pattern {
	return &Pattern{
		Name:    name,
		Value:   value,
		IsRegex: isRegex,
	}
}

func (m *MetaInfo) NewMetaInfo(name string, info []string) *MetaInfo {
	return &MetaInfo{
		Name: name,
		Info: info,
	}
}

func (s *Signature) NewSignature(name string, patterns []Pattern, tags []string, meta []MetaInfo) *Signature {
	return &Signature{
		Name:     name,
		Patterns: patterns,
		Tags:     tags,
		Meta:     meta,
		Matcher:  nil,
	}
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) LoadSignature(data []byte) error {
	var sign Signature
	err := json.Unmarshal(data, &sign) //Десериализуем json

	//Строим массив знаений паттернов и инициализируем ahocorasik.Matcher для него
	patternsValues := []string{}
	for _, pattern := range sign.Patterns {
		patternsValues = append(patternsValues, pattern.Value)
	}

	sign.Matcher = ahocorasick.NewStringMatcher(patternsValues)
	c.Signatures = append(c.Signatures, sign)
	if err != nil {
		return err
	}
	return nil
}

func (c *Compiler) LoadSignatureFromJson(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	c.LoadSignature(data)
	return nil
}

func Match(compiler *Compiler, data string) []*Signature {
	var result []*Signature
	for signIndex := range compiler.Signatures {
		for patternIndex := range compiler.Signatures[signIndex].Patterns {
			if strings.Contains(data, compiler.Signatures[signIndex].Patterns[patternIndex].Value) {
				result = append(result, &compiler.Signatures[signIndex])
			}
		}
	}

	return result
}

func MatchAho(compiler *Compiler, data string) []*Signature {
	var result []*Signature
	for signIndex := range compiler.Signatures {
		if compiler.Signatures[signIndex].Matcher.Match([]byte(data)) != nil {
			result = append(result, &compiler.Signatures[signIndex])
		}
	}

	return result
}

func signsByTagAho(signatures []*Signature, data string) []*Signature {
	var result []*Signature
	for signIndex := range signatures {
		if signatures[signIndex].Matcher.Match([]byte(data)) != nil {
			result = append(result, signatures[signIndex])
		}
	}

	return result
}

func MatchTagsAho(compiler *Compiler, data string, tags []string) []*Signature {
	var taggedSignatures []*Signature //Массив сигнатур с указанными тегами

	//Формируем массив сигнатур с указанными тегами
	for signIndex := range compiler.Signatures {
		for _, tag1 := range tags {
			for _, tag2 := range compiler.Signatures[signIndex].Tags {
				if tag1 == tag2 {
					taggedSignatures = append(taggedSignatures, &compiler.Signatures[signIndex])
					break
				}
			}
		}
	}

	return signsByTagAho(taggedSignatures, data)
}

func signsByTag(signatures []*Signature, data string) []*Signature {
	var result []*Signature
	for signIndex := range signatures {
		for patternIndex := range signatures[signIndex].Patterns {
			if strings.Contains(data, signatures[signIndex].Patterns[patternIndex].Value) {
				result = append(result, signatures[signIndex])
			}
		}
	}

	return result
}

func MatchTags(compiler *Compiler, data string, tags []string) []*Signature {
	var taggedSignatures []*Signature //Массив сигнатур с указанными тегами

	//Формируем массив сигнатур с указанными тегами
	for signIndex := range compiler.Signatures {
		for _, tag1 := range tags {
			for _, tag2 := range compiler.Signatures[signIndex].Tags {
				if tag1 == tag2 {
					taggedSignatures = append(taggedSignatures, &compiler.Signatures[signIndex])
					break
				}
			}
		}
	}

	return signsByTag(taggedSignatures, data)
}
