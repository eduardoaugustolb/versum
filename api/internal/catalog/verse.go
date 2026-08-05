package catalog

type Verse struct {
	BookID  string `json:"book_id"`
	Chapter int    `json:"chapter"`
	Number  int    `json:"number"`
	Text    string `json:"text"`
	Part    int    `json:"part"`
}
