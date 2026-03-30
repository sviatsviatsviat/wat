package shellparse

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestBashParser_andChain(t *testing.T) {
	p := &BashParser{}
	res, err := p.Parse(`cd /tmp/example && git status`)
	if err != nil {
		t.Fatal(err)
	}
	if res.Dialect != core.DialectBash {
		t.Fatalf("dialect %q", res.Dialect)
	}
	if len(res.Pipeline) < 2 {
		t.Fatalf("pipeline: %+v", res.Pipeline)
	}
	if res.Pipeline[0].Name != "cd" {
		t.Errorf("first command: %+v", res.Pipeline[0])
	}
	if res.Pipeline[1].Name != "git" {
		t.Errorf("second command: %+v", res.Pipeline[1])
	}
}

func TestBashParser_pipe(t *testing.T) {
	p := &BashParser{}
	res, err := p.Parse(`echo a | cat`)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Pipeline) != 2 {
		t.Fatalf("want 2 stages, got %+v", res.Pipeline)
	}
	if res.Pipeline[0].PipeLength != 2 || res.Pipeline[1].PipeLength != 2 {
		t.Errorf("pipe lengths: %+v, %+v", res.Pipeline[0], res.Pipeline[1])
	}
}

func TestBashParser_invalidSyntax(t *testing.T) {
	p := &BashParser{}
	_, err := p.Parse(`echo 'unclosed`)
	if err == nil {
		t.Fatal("expected error")
	}
}
