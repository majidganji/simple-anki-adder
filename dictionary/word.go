package dictionary

type Sentence struct {
	Text      string
	Translate string
}

type World struct {
	World             string `form:"world"`
	Audio             string
	Symantec          string
	Type              string
	DefinitionEnglish []string
	DefinitionPersian []string
	Sentences         []Sentence
}
