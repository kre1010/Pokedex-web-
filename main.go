package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"pokedex/api"
	"pokedex/controller"
	"pokedex/models"
)

type PokemonName struct {
	Fr string `json:"fr"`
}

type PokemonSprites struct {
	Regular string `json:"regular"`
}

type PokemonType struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PokemonItem struct {
	Id       int            `json:"pokedex_id"`
	Category string         `json:"category"`
	Name     PokemonName    `json:"name"`
	Sprites  PokemonSprites `json:"sprites"`
	Types    []PokemonType  `json:"types"`
}

func chargerPokemonJSON(chemin string) ([]PokemonItem, error) {
	contenu, err := os.ReadFile(chemin)
	if err != nil {
		return nil, err
	}

	var liste []PokemonItem
	if err := json.Unmarshal(contenu, &liste); err == nil {
		return liste, nil
	}

	var unSeul PokemonItem
	if err := json.Unmarshal(contenu, &unSeul); err == nil && unSeul.Name.Fr != "" {
		return []PokemonItem{unSeul}, nil
	}

	return nil, err
}

func recupererTousLesPokemon() (any, error) {
	liste, err := api.GetAllPokemon()
	if err == nil {
		return liste, nil
	}

	secours, errJSON := chargerPokemonJSON("pokemon.json")
	if errJSON == nil {
		return secours, nil
	}

	return nil, err
}

func chercherPokemonParID(id int) (any, error) {
	p, err := api.GetPokemon(id)
	if err == nil {
		return p, nil
	}

	data, errFallback := recupererTousLesPokemon()
	if errFallback != nil {
		return nil, err
	}

	switch liste := data.(type) {
	case []PokemonItem:
		for _, item := range liste {
			if item.Id == id {
				return item, nil
			}
		}
	case []models.Pokemon:
		return liste, nil
	}

	return nil, err
}

func filtrerParNom(liste []models.Pokemon, recherche string) []models.Pokemon {
	motCle := strings.ToLower(strings.TrimSpace(recherche))
	if motCle == "" {
		return liste
	}

	var resultat []models.Pokemon
	for _, p := range liste {
		nomFR := strings.ToLower(strings.TrimSpace(p.Name.Fr))
		if strings.Contains(nomFR, motCle) {
			resultat = append(resultat, p)
		}
	}

	return resultat
}

func filtrerItemsParNom(items []PokemonItem, recherche string) []PokemonItem {
	motCle := strings.ToLower(strings.TrimSpace(recherche))
	if motCle == "" {
		return items
	}

	var filtres []PokemonItem
	for _, p := range items {
		if strings.Contains(strings.ToLower(strings.TrimSpace(p.Name.Fr)), motCle) {
			filtres = append(filtres, p)
		}
	}

	return filtres
}

func main() {
	tmpl, err := template.ParseGlob("route/templates/*.html")
	if err != nil {
		log.Fatal("Erreur chargement templates : ", err)
	}

	getByIDOrAll := func(idStr string) (any, error) {
		if idStr == "" {
			return recupererTousLesPokemon()
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return nil, models.ErrBadRequest
		}

		return chercherPokemonParID(id)
	}

	ctrl := controller.New(tmpl, getByIDOrAll)

	ctrl.SetSearchByName(func(nom string) (any, error) {
		data, err := recupererTousLesPokemon()
		if err != nil {
			return nil, err
		}

		if liste, ok := data.([]models.Pokemon); ok {
			return filtrerParNom(liste, nom), nil
		}

		if items, ok := data.([]PokemonItem); ok {
			return filtrerItemsParNom(items, nom), nil
		}

		return data, nil
	})

	fs := http.FileServer(http.Dir("route/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", ctrl.Home)
	http.HandleFunc("/pokemon", ctrl.ListPokemon)
	http.HandleFunc("/pokemon/{id}", ctrl.DetailPokemon)

	addr := ":8090"
	fmt.Printf("Serveur lancé sur http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}