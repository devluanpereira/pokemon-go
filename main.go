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

// Estrutura para listar todos os Pokémon
type PokemonList struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Função para buscar os dados do Pokémon por nome
func fetchPokemon(name string) (*Pokemon, error) {
	escapedName := url.QueryEscape(strings.ToLower(strings.TrimSpace(name)))
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", escapedName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar à PokeAPI: %v", err)
	}
	defer resp.Body.Close()

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

// Função para buscar a lista inicial de Pokémon
func fetchPokemonList() ([]Pokemon, error) {
	// Definindo o número de Pokémon que queremos exibir inicialmente
	const limit = 20
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?limit=%d", limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar à PokeAPI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro: PokeAPI retornou status %d", resp.StatusCode)
	}

	var list PokemonList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("erro ao decodificar lista da API: %v", err)
	}

	var pokemons []Pokemon
	for _, result := range list.Results {
		pokemon, err := fetchPokemon(result.Name)
		if err == nil {
			pokemons = append(pokemons, *pokemon)
		}
	}

	return pokemons, nil
}

// Função para renderizar a página inicial com a lista de Pokémon
func indexHandler(w http.ResponseWriter, r *http.Request) {
	pokemons, err := fetchPokemonList()
	if err != nil {
		log.Printf("erro ao buscar lista de Pokémon: %v", err)
		http.Error(w, "Erro ao carregar lista de Pokémon.", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, pokemons); err != nil {
		log.Printf("erro ao renderizar index.html: %v", err)
	}
}

// Função para renderizar o card do Pokémon pesquisado
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
