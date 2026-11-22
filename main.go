package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"pokedex/api"
	"pokedex/controller"
	"pokedex/models"
)

func main() {
	// 1. Charger tous les templates une seule fois au démarrage (dont error.html)
	tmpl, err := template.ParseGlob("route/templates/*.html")
	if err != nil {
		log.Fatal("Erreur chargement templates : ", err)
	}

	// 2. Pont vers api pour gérer liste complète, 400, 404 et 503
	getByIDOrAll := func(idStr string) (any, error) {
		if idStr == "" {
			return api.GetAllPokemon()
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return nil, models.ErrBadRequest
		}

		return api.GetPokemon(id)
	}

	// 3. Créer le contrôleur
	ctrl := controller.New(tmpl, getByIDOrAll)

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