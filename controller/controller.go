package controller

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"pokedex/models"
)

type ErrorData struct {
	Code    int
	Message string
}

type PokemonListPageData struct {
	Pokemons     any
	NomRecherche string
}

type Controller struct {
	templates    *template.Template
	getByIDOrAll func(idStr string) (any, error)
	searchByName func(nom string) (any, error)
}

func New(templates *template.Template, getByIDOrAll func(idStr string) (any, error)) *Controller {
	return &Controller{
		templates:    templates,
		getByIDOrAll: getByIDOrAll,
	}
}

// Permet de brancher la recherche par nom sans casser New(...)
func (c *Controller) SetSearchByName(fn func(nom string) (any, error)) {
	c.searchByName = fn
}

func (c *Controller) showTemplate(w http.ResponseWriter, name string, data any) {
	if err := c.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Println("Erreur template :", err)
	}
}

func (c *Controller) renderError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		log.Println("Erreur :", err)
	}
	w.WriteHeader(status)
	c.showTemplate(w, "error.html", ErrorData{
		Code:    status,
		Message: message,
	})
}

func (c *Controller) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrBadRequest):
		c.renderError(w, http.StatusBadRequest, "Numéro de Pokémon invalide.", err)
	case errors.Is(err, models.ErrNotFound):
		c.renderError(w, http.StatusNotFound, "Ce Pokémon n'existe pas.", err)
	case errors.Is(err, models.ErrUnavailable):
		c.renderError(w, http.StatusServiceUnavailable, "Service momentanément indisponible.", err)
	default:
		c.renderError(w, http.StatusInternalServerError, "Une erreur interne est survenue.", err)
	}
}

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.renderError(w, http.StatusNotFound, "Cette page n'existe pas.", nil)
		return
	}
	http.Redirect(w, r, "/pokemon", http.StatusMovedPermanently)
}

func (c *Controller) ListPokemon(w http.ResponseWriter, r *http.Request) {
	nom := r.URL.Query().Get("nom")

	var (
		pokemons any
		err      error
	)

	if nom != "" && c.searchByName != nil {
		pokemons, err = c.searchByName(nom)
	} else {
		pokemons, err = c.getByIDOrAll("")
	}

	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	c.showTemplate(w, "pokemon.html", PokemonListPageData{
		Pokemons:     pokemons,
		NomRecherche: nom,
	})
}

func (c *Controller) DetailPokemon(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	pokemon, err := c.getByIDOrAll(idStr)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}
	c.showTemplate(w, "detail.html", pokemon)
}