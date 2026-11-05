package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const pikachuJSON = `{"name":"pikachu","weight":60,"types":[{"type":{"name":"electric"}}]}`

// fakeAPI starts a test server, points baseURL at it, and restores state on cleanup.
func fakeAPI(t *testing.T, h http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(h)
	oldBase, oldClient := baseURL, client
	baseURL = srv.URL
	t.Cleanup(func() {
		srv.Close()
		baseURL, client = oldBase, oldClient
	})
}

func TestFetchPokemonSuccess(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pokemon/pikachu" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Write([]byte(pikachuJSON))
	})
	got, err := FetchPokemon("pikachu")
	if err != nil {
		t.Fatal(err)
	}
	want := Pokemon{Name: "pikachu", WeightKg: 6.0, Types: []string{"electric"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFetchPokemonNotFound(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })
	_, err := FetchPokemon("nope")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestFetchPokemonServerError(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) })
	_, err := FetchPokemon("pikachu")
	if err == nil || errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want non-NotFound error", err)
	}
}

func TestFetchPokemonBadJSON(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("garbage")) })
	if _, err := FetchPokemon("pikachu"); err == nil {
		t.Error("expected decode error")
	}
}

func TestFetchPokemonEscapesName(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/pokemon/a%2Fb" {
			t.Errorf("escaped path = %q", r.URL.EscapedPath())
		}
		http.NotFound(w, r)
	})
	FetchPokemon("a/b")
}

func TestFetchPokemonTimeout(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) { time.Sleep(200 * time.Millisecond) })
	client = &http.Client{Timeout: 20 * time.Millisecond}
	if _, err := FetchPokemon("pikachu"); err == nil {
		t.Error("expected timeout error")
	}
}
