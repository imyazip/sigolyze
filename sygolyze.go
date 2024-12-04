package sigolyze

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
