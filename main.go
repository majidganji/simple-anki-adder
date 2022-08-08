package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/majidganji/simple-anki-adder/dictionary"
	"github.com/majidganji/simple-anki-adder/foosoft"
)

type SaveRequest struct {
	Deck  string `form:"deck" json:"deck" binding:"required"`
	World string `form:"world" json:"world" binding:"required"`
	Tags  string `form:"tags" json:"tags"`
}

var worlds map[string]dictionary.World
var fsoft *foosoft.Foosoft

var BackCard = `
<div style="text-align: left;">
{{ .world.Type }}
<hr>
{{ range $i := .world.DefinitionEnglish}}
	- <span>{{ $i }}</span>
	<p></p>
{{ end }}
<hr>
{{ range $i := .world.DefinitionPersian}}
	- <span>{{ $i }}</span><p></p>
{{ end }}
<hr>
{{ range $s := .world.Sentences }}
	<div>
		<p>{{ $s.Text }}</p>
		<p>{{ $s.Translate }}</p>
	</div>
	<hr>		
{{ end }}
</div>
`

func init() {
	worlds = make(map[string]dictionary.World)
	fsoft = foosoft.NewFoosoft()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", Home)
	r.GET("/translate", Translate)
	r.POST("/save", Save)
	r.GET("/refresh-deck-list", RefreshDeckList)

	r.Run(":8000")
}

func Home(c *gin.Context) {
	session := sessions.Default(c)
	messages := session.Flashes()
	search := session.Get("world")
	deck := session.Get("defualt-deck")
	session.Save()

	decks, err := fsoft.GetDecks(false)
	if err != nil {
		panic("cannot get decks from FooSoft Server" + err.Error())
	}

	var world dictionary.World

	if search != nil {
		world = worlds[search.(string)]
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"messages":    messages,
		"world":       world,
		"decks":       decks,
		"worlds":      worlds,
		"defaultDeck": deck,
		"ldoceonline": fmt.Sprintf("https://www.ldoceonline.com/dictionary/%s", strings.ReplaceAll(world.World, " ", "-")),
	})
}

func Translate(c *gin.Context) {
	session := sessions.Default(c)

	world := dictionary.World{
		World: strings.ToLower(c.Query("world")),
	}

	if w, ok := worlds[world.World]; ok && len(w.DefinitionEnglish) != 0 {
		session.Set("world", world.World)
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())

		return
	}

	if world.World == "" {
		session.AddFlash("The world is required.")
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
	}

	df := dictionary.NewDictionaryFacade(&world)
	df.Ldoceonline.Translate()
	df.Fastic.Translate()
	fmt.Println(world)
	if len(world.DefinitionEnglish) == 0 && len(world.DefinitionPersian) == 0 {
		session.AddFlash("The world is wrong.")
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
		return
	}

	worlds[world.World] = world

	session.Set("world", world.World)
	session.Save()

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
}

func Save(c *gin.Context) {
	session := sessions.Default(c)

	var request SaveRequest

	if err := c.ShouldBind(&request); err != nil {
		session.AddFlash(err.Error())
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	world, ok := worlds[request.World]
	if !ok {
		session.AddFlash("The world is not found.")
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	var medias []foosoft.Media
	medias = append(medias, foosoft.Media{
		URL:      world.Audio,
		Filename: world.World,
		Fields:   []string{"Front"},
	})

	tmp := template.Must(template.New("simple").Parse(BackCard))
	var tpl bytes.Buffer
	if err := tmp.Execute(&tpl, map[string]interface{}{"world": world}); err != nil {
		session.AddFlash(err.Error())
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}
	backHtml := tpl.String()

	note := foosoft.Note{
		DeckName:  request.Deck,
		ModelName: "Basic",
		Fields: map[string]string{
			"Front": fmt.Sprintf("<p>%s</p><p>/%s/</p>", world.World, world.Symantec),
			"Back":  backHtml,
		},
		Tags:  strings.Split(request.Tags, ","),
		Audio: medias,
	}

	if err := fsoft.AddNote(note); err != nil {
		session.AddFlash(err.Error())
		session.Save()

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	session.AddFlash("Done")
	session.Set("defualt-deck", request.Deck)
	session.Save()

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func RefreshDeckList(c *gin.Context) {
	fsoft.GetDecks(true)
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
