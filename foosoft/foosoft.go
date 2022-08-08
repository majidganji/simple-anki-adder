package foosoft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Foosoft struct {
	URL     string
	Version int
	Decks   []Deck
}

func NewFoosoft() *Foosoft {
	return &Foosoft{
		URL:     "http://localhost:8765/",
		Version: 6,
	}
}

func (f *Foosoft) GetDecks(update bool) ([]Deck, error) {
	if f.Decks != nil && !update {
		return f.Decks, nil
	}

	res, err := f.sendRequest("deckNamesAndIds", nil)
	if err != nil {
		fmt.Println("err: ", err)
	}

	var decks []Deck

	r := res.(map[string]interface{})
	if r["error"] != nil {
		return nil, fmt.Errorf("%s", r["error"])
	}

	for name, id := range r["result"].(map[string]interface{}) {
		var d Deck
		d.Name = name
		d.ID = id.(float64)

		decks = append(decks, d)
	}

	f.Decks = decks

	return decks, nil
}

func (f *Foosoft) AddNote(note Note) error {
	var notes []Note
	notes = append(notes, note)
	params := make(map[string]interface{})
	params["notes"] = notes

	resutl, err := f.sendRequest("addNotes", params)

	fmt.Println(resutl)

	return err
}

func (f *Foosoft) sendRequest(action string, params map[string]interface{}) (interface{}, error) {
	data := make(map[string]interface{})
	data["action"] = action
	data["version"] = f.Version

	if params != nil {
		data["params"] = params
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(f.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}

	json.NewDecoder(response.Body).Decode(&res)

	return res, nil
}
