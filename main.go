package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
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
	ID        int            `json:"id"`
	Category  string         `json:"category"`
	Name      PokemonName    `json:"name"`
	Sprites   PokemonSprites `json:"sprites"`
	Types     []PokemonType  `json:"types"`
}

func (p *PokemonItem) UnmarshalJSON(data []byte) error {
	type Alias PokemonItem
	aux := (*Alias)(p)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	// Si le JSON utilise "id" au lieu de "pokedex_id", on récupère l'ID
	if p.PokedexID == 0 && p.ID != 0 {
		p.PokedexID = p.ID
	}
	return nil
}

func loadPokemonFromJSON(path string) ([]PokemonItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimSpace(data)

	// Cas 1 : le JSON commence par '[' -> tableau direct
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var list []PokemonItem
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return nil, err
		}
		return list, nil
	}

	// Cas 2 : le JSON commence par '{' -> objet avec une clé liste (ex: pokemons, data, results, items...)
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &wrapper); err != nil {
		return nil, err
	}

	for _, key := range []string{"pokemons", "pokemon", "data", "results", "items", "list"} {
		if raw, ok := wrapper[key]; ok {
			var list []PokemonItem
			if err := json.Unmarshal(raw, &list); err == nil && len(list) > 0 {
				return list, nil
			}
		}
	}

	// Si l'objet contient directement des clés de Pokémon (ex: un seul Pokémon ou une map par ID)
	var single PokemonItem
	if err := json.Unmarshal(trimmed, &single); err == nil && (single.PokedexID != 0 || single.Name.Fr != "") {
		return []PokemonItem{single}, nil
	}

	var byKey map[string]PokemonItem
	if err := json.Unmarshal(trimmed, &byKey); err == nil && len(byKey) > 0 {
		var list []PokemonItem
		for _, item := range byKey {
			if item.PokedexID != 0 || item.Name.Fr != "" {
				list = append(list, item)
			}
		}
		if len(list) > 0 {
			return list, nil
		}
	}

	return nil, fmt.Errorf("format JSON non reconnu dans %s", path)
}

func getAllPokemonSafe() (any, error) {
	list, err := api.GetAllPokemon()
	if err == nil {
		return list, nil
	}

	// Secours sur pokemon.json si l'API échoue ou dépasse le délai
	fallback, jsonErr := loadPokemonFromJSON("pokemon.json")
	if jsonErr == nil {
		return fallback, nil
	}

	return nil, err
}

func extractFrenchName(v reflect.Value) string {
	if v.Kind() != reflect.Struct {
		return ""
	}

	if f := v.FieldByName("NomFR"); f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}

	if nameField := v.FieldByName("Name"); nameField.IsValid() {
		if nameField.Kind() == reflect.Ptr && !nameField.IsNil() {
			nameField = nameField.Elem()
		}
		if nameField.Kind() == reflect.Struct {
			for _, fieldName := range []string{"Fr", "FR", "French"} {
				if sub := nameField.FieldByName(fieldName); sub.IsValid() && sub.Kind() == reflect.String {
					return sub.String()
				}
			}
		}
	}

	return ""
}

func filterPokemonByName(all any, nom string) any {
	// 1. Normalisation : supprimer les espaces en trop et mettre en minuscules
	query := strings.ToLower(strings.TrimSpace(nom))
	if query == "" {
		return all
	}

	val := reflect.ValueOf(all)
	if val.Kind() != reflect.Slice {
		return all
	}

	// 2. Filtrage partiel sur le nom français
	filtered := reflect.MakeSlice(val.Type(), 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		target := elem
		if target.Kind() == reflect.Ptr {
			if target.IsNil() {
				continue
			}
			target = target.Elem()
		}

		nomFR := strings.ToLower(strings.TrimSpace(extractFrenchName(target)))
		if strings.Contains(nomFR, query) {
			filtered = reflect.Append(filtered, elem)
		}
	}

	return filtered.Interface()
}

func main() {
	// 1. Charger tous les templates une seule fois au démarrage (dont error.html)
	tmpl, err := template.ParseGlob("route/templates/*.html")
	if err != nil {
		log.Fatal("Erreur chargement templates : ", err)
	}

	// 2. Pont vers api (avec secours pokemon.json) pour gérer liste complète, 400, 404 et 503
	getByIDOrAll := func(idStr string) (any, error) {
		if idStr == "" {
			return getAllPokemonSafe()
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return nil, models.ErrBadRequest
		}

		return api.GetPokemon(id)
	}

	// 3. Créer le contrôleur + brancher la recherche partielle par nom FR
	ctrl := controller.New(tmpl, getByIDOrAll)
	ctrl.SetSearchByName(func(nom string) (any, error) {
		all, err := getAllPokemonSafe()
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