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
		return nil, err
	}
	request.Header.Set("User-Agent", "DataDex/1.0")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("réponse inattendue : %s", response.Status)
	}

	var pokemons []models.Pokemon
	err = json.NewDecoder(response.Body).Decode(&pokemons)
	return pokemons, err
}
