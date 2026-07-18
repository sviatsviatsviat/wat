package hookkit

import "testing"

func TestRawPayload(t *testing.T) {
	t.Parallel()
	var p RawPayload
	if p.Raw() != nil {
		t.Fatal("empty RawPayload should return nil")
	}
	p.SetRaw([]byte(`{"a":1}`))
	got := p.Raw()
	if string(got) != `{"a":1}` {
		t.Fatalf("Raw() = %q", got)
	}
	got[0] = 'x'
	if string(p.Raw()) != `{"a":1}` {
		t.Fatal("Raw() must clone")
	}
}
