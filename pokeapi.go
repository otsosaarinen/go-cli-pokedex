package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrNotFound = errors.New("pokemon not found")

	client = &http.Client{Timeout: 10 * time.Second}

	// baseURL is a variable so tests can point it at a fake server.
	baseURL = "https://pokeapi.co/api/v2"
)

type Pokemon struct {
	Name     string
	WeightKg float64
	Types    []string
}

type apiResponse struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"` // hectograms
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func FetchPokemon(name string) (Pokemon, error) {
	resp, err := client.Get(baseURL + "/pokemon/" + url.PathEscape(name))
	if err != nil {
		return Pokemon{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Pokemon{}, ErrNotFound
	default:
		return Pokemon{}, fmt.Errorf("pokeapi returned %s", resp.Status)
	}

	var raw apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Pokemon{}, fmt.Errorf("decoding response: %w", err)
	}

	p := Pokemon{
		Name:     raw.Name,
		WeightKg: float64(raw.Weight) / 10,
	}
	for _, t := range raw.Types {
		p.Types = append(p.Types, t.Type.Name)
	}
	return p, nil
}
