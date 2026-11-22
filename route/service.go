package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrBadRequest = errors.New("identifiant invalide : doit être un nombre positif")

type TypeInfo struct {
	Nom  string `json:"nom"`
	Icon string `json:"icon"`
}

type PokemonDetail struct {
	ID        int        `json:"id"`
	Numero    int        `json:"numero"`
	NomFR     string     `json:"nom_fr"`
	Categorie string     `json:"categorie"`
	ImageURL  string     `json:"image"`
	Types     []TypeInfo `json:"types"`
}

type PokemonService struct {
	BaseURL string
	Client  *http.Client
}

func NewPokemonService(baseURL string) *PokemonService {
	return &PokemonService{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *PokemonService) GetByIDOrAll(idStr string) ([]PokemonDetail, error) {
	if idStr == "" {
		resp, err := s.Client.Get(s.BaseURL + "/api/pokemon")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var list []PokemonDetail
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			return nil, err
		}
		return list, nil
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return nil, ErrBadRequest
	}

	url := fmt.Sprintf("%s/api/pokemon/%d", s.BaseURL, id)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []PokemonDetail{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API : statut %d", resp.StatusCode)
	}

	var one PokemonDetail
	if err := json.NewDecoder(resp.Body).Decode(&one); err != nil {
		return nil, err
	}

	return []PokemonDetail{one}, nil
}

func (s *PokemonService) SearchByName(nom string) ([]PokemonDetail, error) {
	all, err := s.GetByIDOrAll("")
	if err != nil {
		return nil, err
	}

	query := strings.ToLower(strings.TrimSpace(nom))
	if query == "" {
		return all, nil
	}

	var filtered []PokemonDetail
	for _, p := range all {
		nomFR := strings.ToLower(strings.TrimSpace(p.NomFR))
		if strings.Contains(nomFR, query) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}