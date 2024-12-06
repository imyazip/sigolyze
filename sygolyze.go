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
	Signatures      []Signature
	SignaturesByTag map[string][]*Signature
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
	return &Compiler{
		SignaturesByTag: make(map[string][]*Signature),
	}
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

	for _, tag := range sign.Tags {
		c.SignaturesByTag[tag] = append(c.SignaturesByTag[tag], &sign)
	}
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

func GetSignaturesByTags(compiler *Compiler, tags []string) []*Signature {
	signatureSet := make(map[*Signature]struct{}) // Используем map для уникальности сигнатур

	// Ищем сигнатуры по каждому тегу
	for _, tag := range tags {
		if signatures, exists := compiler.SignaturesByTag[tag]; exists {
			for _, signature := range signatures {
				signatureSet[signature] = struct{}{} // Добавляем в set
			}
		}
	}

	// Преобразуем set в массив
	var result []*Signature
	for signature := range signatureSet {
		result = append(result, signature)
	}

	return result
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
	return signsByTagAho(GetSignaturesByTags(compiler, tags), data)
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
	return signsByTag(GetSignaturesByTags(compiler, tags), data)
}
