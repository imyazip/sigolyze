package sigolyze

import (
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
}

type Compiler struct {
	Signatures Signature
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
	}
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) LoadRules(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &c.Signatures)
	if err != nil {
		return err
	}
	return nil
}
