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

func GetPokemon(id int) (models.Pokemon, error) {
	endpoint := fmt.Sprintf("%s/pokemon/%d", baseURL, id)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return models.Pokemon{}, err
	}
	request.Header.Set("User-Agent", "DataDex/1.0")

	response, err := client.Do(request)
	if err != nil {
		return models.Pokemon{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return models.Pokemon{}, fmt.Errorf("réponse inattendue : %s",
			response.Status)
	}
	var pokemon models.Pokemon
	err = json.NewDecoder(response.Body).Decode(&pokemon)
	return pokemon, err

}
