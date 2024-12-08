package sigolyze

import (
	"fmt"
	"io"
	"os"
	"regexp"

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
	Name            string     `json:"name"`
	Patterns        []Pattern  `json:"patterns"`
	Tags            []string   `json:"tags"`
	Meta            []MetaInfo `json:"meta"`
	Matcher         *ahocorasick.Matcher
	regexpCompilers []*regexp.Regexp
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
	regexPatternsValues := []*regexp.Regexp{}
	stringPatternsValues := []string{}
	for _, pattern := range sign.Patterns {
		if pattern.IsRegex {
			compiled, err := regexp.Compile(pattern.Value)
			if err != nil {
				return fmt.Errorf("failed to parse regex pattern %s in signature %s. Error: %s", pattern.Value, sign.Name, err)
			}
			regexPatternsValues = append(regexPatternsValues, compiled)
		} else {
			stringPatternsValues = append(stringPatternsValues, pattern.Value)
		}
	}
	for _, regexPattern := range regexPatternsValues {
		sign.regexpCompilers = append(sign.regexpCompilers, regexPattern)
	}

	sign.Matcher = ahocorasick.NewStringMatcher(stringPatternsValues)
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
		if compiler.Signatures[signIndex].Matcher.Match([]byte(data)) != nil {
			result = append(result, &compiler.Signatures[signIndex])
		}

		if len(compiler.Signatures[signIndex].regexpCompilers) != 0 {
			for _, regex := range compiler.Signatures[signIndex].regexpCompilers {
				if regex.Match([]byte(data)) {
					result = append(result, &compiler.Signatures[signIndex])
				}
			}
		}
	}

	return result
}

func signsByTag(signatures []*Signature, data string) []*Signature {
	var result []*Signature
	for signIndex := range signatures {
		if signatures[signIndex].Matcher.Match([]byte(data)) != nil {
			result = append(result, signatures[signIndex])
		}

		if len(signatures[signIndex].regexpCompilers) != 0 {
			for _, regex := range signatures[signIndex].regexpCompilers {
				if regex.Match([]byte(data)) {
					result = append(result, signatures[signIndex])
				}
			}
		}
	}

	return result
}

func MatchTags(compiler *Compiler, data string, tags []string) []*Signature {
	return signsByTag(GetSignaturesByTags(compiler, tags), data)
}
