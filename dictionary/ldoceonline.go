package dictionary

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Ldoceonline struct {
	baseurl   string
	Collector *colly.Collector
	World     *World
}

func NewLdoceonline(world *World) *Ldoceonline {

	return &Ldoceonline{
		baseurl:   "https://www.ldoceonline.com",
		Collector: colly.NewCollector(),
		World:     world,
	}
}

func (l *Ldoceonline) Translate() error {
	l.Collector.OnHTML(".Sense > .DEF", func(e *colly.HTMLElement) {
		l.World.DefinitionEnglish = append(l.World.DefinitionEnglish, strings.TrimSpace(e.Text))
	})

	l.Collector.OnHTML(".POS", func(e *colly.HTMLElement) {
		t := strings.TrimSpace(e.Text)
		if !strings.Contains(l.World.Type, t) {
			l.World.Type += fmt.Sprintf("%s, ", t)
		}
	})

	l.Collector.OnHTML("[title=\"Play American pronunciation of "+l.World.World+"\"]", func(e *colly.HTMLElement) {
		l.World.Audio = e.Attr("data-src-mp3")
	})

	l.Collector.Visit(l.baseurl + "/dictionary/" + strings.ReplaceAll(l.World.World, " ", "-"))

	return nil
}
