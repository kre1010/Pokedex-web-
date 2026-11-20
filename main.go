package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pokedex/api"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pokemons, err := api.GetAllPokemon()
	if err != nil {
		http.Error(w, "Erreur récupération Pokémon : "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseGlob("route/templates/*.html")
	if err != nil {
		http.Error(w, "Erreur chargement template : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "pokemon.html", pokemons)
	if err != nil {
		http.Error(w, "Erreur affichage page : "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("route/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	addr := ":8090"
	fmt.Printf("Serveur lancé sur http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}