package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cl *Client) ListPokemon(pageURL *string) (PokemonResp, error) {
	fullURL := baseURL + "/pokemon/?offset=0&limit=20"
	if pageURL != nil {
		fullURL = *pageURL
	}

	// check the cache

	if data, exists := cl.cache.Get(fullURL); exists {
		fmt.Println("Cache hit!")
		pokemonResp := PokemonResp{}

		err := json.Unmarshal(data, &pokemonResp)
		if err != nil {
			return PokemonResp{}, err
		}

		return pokemonResp, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonResp{}, err
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return PokemonResp{}, err
	}

	if resp.StatusCode > 399 {
		return PokemonResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonResp{}, err
	}

	// add to cache
	cl.cache.Add(fullURL, data)

	pokemonResp := PokemonResp{}

	err = json.Unmarshal(data, &pokemonResp)
	if err != nil {
		return PokemonResp{}, err
	}

	return pokemonResp, nil
}

func (cl *Client) ExplorePokemon(specificPokemon *string) (SpecificPokemonResp, error) {

	fullURL := baseURL + "/pokemon/" + *specificPokemon

	// check the cache

	if data, exists := cl.cache.Get(fullURL); exists {
		fmt.Println("Cache hit!")
		specificPokemonResp := SpecificPokemonResp{}

		err := json.Unmarshal(data, &specificPokemonResp)
		if err != nil {
			return SpecificPokemonResp{}, err
		}

		return specificPokemonResp, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return SpecificPokemonResp{}, err
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return SpecificPokemonResp{}, err
	}

	if resp.StatusCode > 399 {
		return SpecificPokemonResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpecificPokemonResp{}, err
	}

	// add to cache
	cl.cache.Add(fullURL, data)

	specificPokemonResp := SpecificPokemonResp{}

	err = json.Unmarshal(data, &specificPokemonResp)
	if err != nil {
		return SpecificPokemonResp{}, err
	}

	return specificPokemonResp, nil
}
