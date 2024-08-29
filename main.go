package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Estrutura para armazenar os dados do Pokémon

type Pokemon struct {
    Name string `json:"name"`
    Sprites struct {
        FrontDefault string `json:"front_default"`
    } `json:"sprites"`
}

func fetchPokemon(name string) (*Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}