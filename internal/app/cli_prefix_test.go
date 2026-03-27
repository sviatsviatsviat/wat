package app

import (
	"reflect"
	"testing"
)

func TestParseHost_minimal(t *testing.T) {
	hookHostName, argsAfterHost, err := parseHost([]string{"cursor", "run", "echo"})
	if err != nil {
		t.Fatal(err)
	}
	if hookHostName != "cursor" || !reflect.DeepEqual(argsAfterHost, []string{"run", "echo"}) {
		t.Fatalf("hookHostName=%q argsAfterHost=%v", hookHostName, argsAfterHost)
	}
}

func TestParseHost_emptyHost(t *testing.T) {
	_, _, err := parseHost([]string{"", "run"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseSubcommand_minimal(t *testing.T) {
	watSubcommand, subcommandArgs, err := parseSubcommand([]string{"run", "echo", "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if watSubcommand != "run" || !reflect.DeepEqual(subcommandArgs, []string{"echo", "hi"}) {
		t.Fatalf("watSubcommand=%q subcommandArgs=%v", watSubcommand, subcommandArgs)
	}
}

func TestParseSubcommand_missing(t *testing.T) {
	_, _, err := parseSubcommand([]string{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseSubcommand_emptyCommand(t *testing.T) {
	_, _, err := parseSubcommand([]string{"", "run"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseHostThenSubcommand_chain(t *testing.T) {
	programArgs := []string{"cursor", "run", "echo", "hi"}
	hookHostName, argsAfterHost, err := parseHost(programArgs)
	if err != nil {
		t.Fatal(err)
	}
	watSubcommand, subcommandArgs, err := parseSubcommand(argsAfterHost)
	if err != nil {
		t.Fatal(err)
	}
	if hookHostName != "cursor" || watSubcommand != "run" || !reflect.DeepEqual(subcommandArgs, []string{"echo", "hi"}) {
		t.Fatalf("hookHostName=%q watSubcommand=%q subcommandArgs=%v", hookHostName, watSubcommand, subcommandArgs)
	}
}
