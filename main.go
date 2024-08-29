package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Estrutura para armazenar os dados do Pokémon

type Pokemon struct {
    Name string `json:"name"`
    Sprites struct {
        FrontDefault string `json:"front_default"`
    } `json:"sprites"`
}

// Função para buscar os dados do Pokemon da PokeApi

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

// Função para renderizar a pagina inicial

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// Função para rendenrizar o card do pokemon

func pokemonHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Faltou o nome do Pokemon.", http.StatusBadRequest)
		return
	}

	pokemon, err := fetchPokemon(name)
	if err != nil {
		http.Error(w, "Erro ao buscar o Pokemon.", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/pokemon.html"))
	tmpl.Execute(w, pokemon)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pokemon", pokemonHandler)
	fmt.Println("Servidor rodando na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}