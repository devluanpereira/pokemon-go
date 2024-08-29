package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

// Estrutura para armazenar os dados do Pokémon
type Pokemon struct {
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

// Função para buscar os dados do Pokémon da PokeAPI
func fetchPokemon(name string) (*Pokemon, error) {
	// Sanitizando o nome do Pokémon
	escapedName := url.QueryEscape(strings.ToLower(strings.TrimSpace(name)))
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", escapedName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar à PokeAPI: %v", err)
	}
	defer resp.Body.Close()

	// Verificando se o status da resposta é 200 OK
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Pokémon não encontrado. Verifique o nome e tente novamente.")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro: PokeAPI retornou status %d", resp.StatusCode)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta da API: %v", err)
	}

	return &pokemon, nil
}

// Função para renderizar a página inicial
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("erro ao renderizar index.html: %v", err)
	}
}

// Função para renderizar o card do Pokémon
func pokemonHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Faltou o nome do Pokémon.", http.StatusBadRequest)
		return
	}

	pokemon, err := fetchPokemon(name)
	if err != nil {
		log.Printf("erro ao buscar o Pokémon: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/pokemon.html"))
	if err := tmpl.Execute(w, pokemon); err != nil {
		log.Printf("erro ao renderizar pokemon.html: %v", err)
	}
}

func main() {
	// Handler para arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pokemon", pokemonHandler)
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
