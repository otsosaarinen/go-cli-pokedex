package main

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestRunMultiple(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/pikachu") {
			w.Write([]byte(pikachuJSON))
			return
		}
		http.NotFound(w, r)
	})

	var out, errOut bytes.Buffer
	code := run([]string{"get", "Pikachu", "missing", "pikachu"}, &out, &errOut)

	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	block := "Name: pikachu\nWeight: 6.0 kg\nTypes: electric\n"
	if want := block + "\n" + block; out.String() != want {
		t.Errorf("stdout = %q, want %q", out.String(), want)
	}
	if want := "\nerror: missing: pokemon not found\n"; errOut.String() != want {
		t.Errorf("stderr = %q, want %q", errOut.String(), want)
	}
}

func TestRunSuccessExitsZero(t *testing.T) {
	fakeAPI(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(pikachuJSON)) })
	var out, errOut bytes.Buffer
	if code := run([]string{"get", "pikachu"}, &out, &errOut); code != 0 {
		t.Errorf("code = %d, want 0", code)
	}
}

func TestRunUsage(t *testing.T) {
	for _, args := range [][]string{nil, {"get"}, {"fetch", "pikachu"}} {
		var out, errOut bytes.Buffer
		if code := run(args, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "usage:") {
			t.Errorf("args %v: code=%d stderr=%q", args, code, errOut.String())
		}
	}
}
