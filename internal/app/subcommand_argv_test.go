package app

import (
	"reflect"
	"strings"
	"testing"
)

func TestConsumeSubcommandHostFlags_None(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"echo", "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "hi"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_One(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"--host", "cursor", "go", "version"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"go", "version"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_DashCapitalH(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"-H", "cursor", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || !reflect.DeepEqual(subcommandArgs, []string{"x"}) {
		t.Fatalf("host=%q subcommandArgs=%#v", host, subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_LastWins(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"--host", "a", "-H", "cursor", "run"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"run"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_MissingValue(t *testing.T) {
	_, _, err := consumeSubcommandHostFlags([]string{"--host"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--host") {
		t.Fatalf("error %q should mention --host", err.Error())
	}
}

func TestConsumeSubcommandHostFlags_MissingValueDashCapitalH(t *testing.T) {
	_, _, err := consumeSubcommandHostFlags([]string{"-H"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "-H") {
		t.Fatalf("error %q should mention -H", err.Error())
	}
}

func TestConsumeSubcommandHostFlags_EqualsFormLong(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"--host=cursor", "run", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || !reflect.DeepEqual(subcommandArgs, []string{"run", "x"}) {
		t.Fatalf("host=%q subcommandArgs=%#v", host, subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_DashHEqualsNotParsedAsHost(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"-H=cursor", "y"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want default cursor", host)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"-H=cursor", "y"}) {
		t.Fatalf("subcommandArgs %#v, want -H= left intact", subcommandArgs)
	}
}

func TestConsumeSubcommandHostFlags_EqualsFormEmptyValue(t *testing.T) {
	_, _, err := consumeSubcommandHostFlags([]string{"--host=", "run"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--host") {
		t.Fatalf("error %q should mention --host", err.Error())
	}
}

func TestConsumeSubcommandHostFlags_EqualsFormLastWins(t *testing.T) {
	host, subcommandArgs, err := consumeSubcommandHostFlags([]string{"--host=a", "-H", "cursor", "run"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || !reflect.DeepEqual(subcommandArgs, []string{"run"}) {
		t.Fatalf("host=%q subcommandArgs=%#v", host, subcommandArgs)
	}
}
