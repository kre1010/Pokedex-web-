package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokedex/models"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

const baseURL = "https://tyradex.app/api/v1"

func GetAllPokemon() ([]models.Pokemon, error) {
	endpoint := fmt.Sprintf("%s/pokemon", baseURL)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrUnavailable, err)
	}
	request.Header.Set("User-Agent", "DataDex/1.0")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrUnavailable, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, models.ErrNotFound
	}
	if response.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("%w: statut %s", models.ErrUnavailable, response.Status)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("réponse inattendue : %s", response.Status)
	}

	var pokemons []models.Pokemon
	if err := json.NewDecoder(response.Body).Decode(&pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

func GetPokemon(id int) (models.Pokemon, error) {
	endpoint := fmt.Sprintf("%s/pokemon/%d", baseURL, id)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return models.Pokemon{}, fmt.Errorf("%w: %v", models.ErrUnavailable, err)
	}
	request.Header.Set("User-Agent", "DataDex/1.0")

	response, err := client.Do(request)
	if err != nil {
		return models.Pokemon{}, fmt.Errorf("%w: %v", models.ErrUnavailable, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return models.Pokemon{}, models.ErrNotFound
	}
	if response.StatusCode >= http.StatusInternalServerError {
		return models.Pokemon{}, fmt.Errorf("%w: statut %s", models.ErrUnavailable, response.Status)
	}
	if response.StatusCode != http.StatusOK {
		return models.Pokemon{}, fmt.Errorf("réponse inattendue : %s", response.Status)
	}

	var pokemon models.Pokemon
	if err := json.NewDecoder(response.Body).Decode(&pokemon); err != nil {
		return models.Pokemon{}, models.ErrNotFound
	}

	if len(pokemon.Types) == 0 {
		return models.Pokemon{}, models.ErrNotFound
	}

	return pokemon, nil
}