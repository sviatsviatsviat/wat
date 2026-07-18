package hookkit

import (
	"encoding/json"
	"testing"
)

func TestDecodeAsAndThen(t *testing.T) {
	t.Parallel()
	type ev struct {
		Name  string `json:"name"`
		Bound string
	}
	got, err := DecodeAsAndThen(json.RawMessage(`{"name":"x"}`), func(e *ev, _ []byte) {
		e.Bound = e.Name
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "x" || got.Bound != "x" {
		t.Fatalf("got = %+v", got)
	}
	_, err = DecodeAsAndThen[ev]([]byte(`{`), nil)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestDecodeAsAndThen_SetsRawPayload(t *testing.T) {
	t.Parallel()
	type ev struct {
		RawPayload
		Name string `json:"name"`
	}
	raw := []byte(`{"name":"x"}`)
	got, err := DecodeAsAndThen[ev](raw, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(got.Raw(), raw) {
		t.Fatalf("Raw() = %s, want %s", got.Raw(), raw)
	}
}

func jsonEqual(a, b []byte) bool {
	return string(a) == string(b)
}
