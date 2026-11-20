package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

type DetailPageData struct {
	Query    string
	Pokemons []PokemonDetail
	IsList   bool
	ErrorMsg string
}

func main() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	svc := NewPokemonService("http://localhost:8080")

	mux := http.NewServeMux()

	// Route GET /pokemon/{id}
	mux.HandleFunc("GET /pokemon/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")

		pokemons, err := svc.GetByIDOrAll(idParam)
		if err != nil {
			// Si le service renvoie ErrBadRequest, on affiche un message clair (400)
			if errors.Is(err, ErrBadRequest) {
				w.WriteHeader(http.StatusBadRequest)
				data := DetailPageData{
					Query:    idParam,
					ErrorMsg: "Aucun Pokémon ne correspond à " + idParam,
				}
				_ = tmpl.ExecuteTemplate(w, "detail.html", data)
				return
			}

			http.Error(w, "Erreur serveur : "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := DetailPageData{
			Query:    idParam,
			Pokemons: pokemons,
			IsList:   idParam == "",
		}
		_ = tmpl.ExecuteTemplate(w, "detail.html", data)
	})

	// Si l'URL est /pokemon/ avec rien après {id}, on affiche aussi la liste complète
	mux.HandleFunc("GET /pokemon/", func(w http.ResponseWriter, r *http.Request) {
		pokemons, err := svc.GetByIDOrAll("")
		if err != nil {
			http.Error(w, "Erreur serveur : "+err.Error(), http.StatusInternalServerError)
			return
		}
		data := DetailPageData{
			Pokemons: pokemons,
			IsList:   true,
		}
		_ = tmpl.ExecuteTemplate(w, "detail.html", data)
	})

	log.Println("Serveur lancé sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
