package dictionary

import (
	"strings"

	"github.com/gocolly/colly"
)

type Fastic struct {
	baseurl   string
	Collector *colly.Collector
	World     *World
}

func NewFastic(world *World) *Fastic {
	return &Fastic{
		baseurl:   "https://fastdic.com",
		Collector: colly.NewCollector(),
		World:     world,
	}
}

func (f *Fastic) Translate() error {
	f.Collector.OnHTML(".result > li", func(e *colly.HTMLElement) {
		text := e.DOM.Contents().Nodes[2].Data
		text = strings.TrimSpace(text)
		texts := strings.Split(text, "،")
		f.World.DefinitionPersian = append(f.World.DefinitionPersian, texts...)
	})

	f.Collector.OnHTML(".results__phonetics > li:nth-child(1) > strong:nth-child(2)", func(e *colly.HTMLElement) {
		f.World.Symantec = strings.TrimSpace(e.Text)
	})

	f.Collector.OnHTML(".result > li > div > ul", func(e *colly.HTMLElement) {
		var sentence Sentence
		sentence.Text = strings.TrimSpace(e.DOM.Children().First().Text())
		sentence.Text = strings.TrimLeft(sentence.Text, "- ")

		sentence.Translate = strings.TrimSpace(e.DOM.Children().Last().Text())
		sentence.Translate = strings.TrimLeft(sentence.Translate, "- ")

		f.World.Sentences = append(f.World.Sentences, sentence)
	})

	f.Collector.Visit(f.baseurl + "/word/" + f.World.World)

	return nil
}
