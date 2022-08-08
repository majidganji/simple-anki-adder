package foosoft

type Media struct {
	URL      string   `json:"url"`
	Filename string   `json:"filename"`
	SkipHash string   `json:"skipHash"`
	Fields   []string `json:"fields"`
}

type Note struct {
	DeckName  string            `json:"deckName"`
	ModelName string            `json:"modelName"`
	Fields    map[string]string `json:"fields"`
	Tags      []string          `json:"tags"`
	Audio     []Media           `json:"audio"`
	Video     []Media           `json:"video"`
	Picture   []Media           `json:"picture"`
}
