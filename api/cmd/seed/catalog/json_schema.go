package main

type CorpusManifest struct {
	BibleSha256 string `json:"bibleSha256"`
}

type CorpusBible struct {
	SchemaVersion int    `json:"schemaVersion"`
	Books         []Book `json:"books"`
}

type Book struct {
	Id        string    `json:"id"`
	Order     int       `json:"order"`
	Name      string    `json:"name"`
	Testament string    `json:"testament"`
	Chapters  []Chapter `json:"chapters"`
}

type Chapter struct {
	Number int     `json:"number"`
	Verses []Verse `json:"verses"`
}

type Verse struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
	Part   int    `json:"part"`
}
