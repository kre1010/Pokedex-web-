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
	PokedexID int            `json:"pokedex_id"`
	Category  string         `json:"category"`
	Name      PokemonName    `json:"name"`
	Sprites   PokemonSprites `json:"sprites"`
	Types     []PokemonType  `json:"types"`
}

func loadPokemonFromJSON(path string) ([]PokemonItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list []PokemonItem
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func filterPokemonByName(list []PokemonItem, nom string) []PokemonItem {
	// 1. Normalisation : supprimer les espaces en trop et mettre en minuscules
	query := strings.ToLower(strings.TrimSpace(nom))
	if query == "" {
		return list
	}

	// 2. Filtrage partiel sur le nom français
	var result []PokemonItem
	for _, p := range list {
		nomFR := strings.ToLower(strings.TrimSpace(p.Name.Fr))
		if strings.Contains(nomFR, query) {
			result = append(result, p)
		}
	}

	return result
}

func main() {
	// 1. Charger tous les templates une seule fois au démarrage (dont error.html)
	tmpl, err := template.ParseGlob("route/templates/*.html")
	if err != nil {
		log.Fatal("Erreur chargement templates : ", err)
	}

	// 2. Pont vers api pour gérer détail / 400 / 404 / 503
	getByIDOrAll := func(idStr string) (any, error) {
		if idStr == "" {
			return loadPokemonFromJSON("pokemon.json")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return nil, models.ErrBadRequest
		}

		return api.GetPokemon(id)
	}

	// 3. Créer le contrôleur + brancher la recherche par nom FR depuis pokemon.json
	ctrl := controller.New(tmpl, getByIDOrAll)
	ctrl.SetSearchByName(func(nom string) (any, error) {
		all, err := loadPokemonFromJSON("pokemon.json")
		if err != nil {
			return nil, err
		}
		return filterPokemonByName(all, nom), nil
	})

	// 4. Routes statiques et applicatives
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