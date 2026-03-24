package app

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestConsumeSubcommandSharedFlags_None(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"echo", "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if filePattern != nil {
		t.Fatalf("filePattern %#v, want nil", filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "hi"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_One(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--host", "cursor", "go", "version"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if filePattern != nil {
		t.Fatalf("filePattern %#v, want nil", filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"go", "version"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_DashCapitalH(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"-H", "cursor", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern != nil || !reflect.DeepEqual(subcommandArgs, []string{"x"}) {
		t.Fatalf("host=%q filePattern=%v subcommandArgs=%#v", host, filePattern, subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_LastWins(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--host", "a", "-H", "cursor", "run"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" {
		t.Fatalf("host %q, want cursor", host)
	}
	if filePattern != nil {
		t.Fatalf("filePattern %#v, want nil", filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"run"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_MissingValue(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"--host"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "needs an argument") || !strings.Contains(err.Error(), "-host") {
		t.Fatalf("error %q should be flag needs-an-argument for host", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_MissingValueDashCapitalH(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"-H"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "-H") {
		t.Fatalf("error %q should mention -H", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_EqualsFormLong(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--host=cursor", "run", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern != nil || !reflect.DeepEqual(subcommandArgs, []string{"run", "x"}) {
		t.Fatalf("host=%q filePattern=%v subcommandArgs=%#v", host, filePattern, subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_DashHEqualsHost(t *testing.T) {
	// Standard [flag] parsing treats -H=value like other string flags (same idea as --host=value).
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"-H=cursor", "y"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern != nil || !reflect.DeepEqual(subcommandArgs, []string{"y"}) {
		t.Fatalf("host=%q filePattern=%v subcommandArgs=%#v", host, filePattern, subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_EqualsFormEmptyValue(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"--host=", "run"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--host") {
		t.Fatalf("error %q should mention --host", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_EqualsFormLastWins(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--host=a", "-H", "cursor", "run"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern != nil || !reflect.DeepEqual(subcommandArgs, []string{"run"}) {
		t.Fatalf("host=%q filePattern=%v subcommandArgs=%#v", host, filePattern, subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternLongEquals(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--file-pattern=[.]go$", "echo", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern == nil || *filePattern != `[.]go$` {
		t.Fatalf("host=%q filePattern=%v", host, filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "x"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternLongSeparate(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"--file-pattern", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern == nil || *filePattern != `[.]go$` {
		t.Fatalf("host=%q filePattern=%v", host, filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "x"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternShort(t *testing.T) {
	_, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if filePattern == nil || *filePattern != `[.]go$` {
		t.Fatalf("filePattern=%v", filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "x"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternLastWins(t *testing.T) {
	_, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"-f", "a", "--file-pattern=b", "echo"})
	if err != nil {
		t.Fatal(err)
	}
	if filePattern == nil || *filePattern != "b" {
		t.Fatalf("filePattern=%v want b", filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_HostAndPatternInterleaved(t *testing.T) {
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags([]string{"-f", `[.]go$`, "-H", "cursor", "echo", "x"})
	if err != nil {
		t.Fatal(err)
	}
	if host != "cursor" || filePattern == nil || *filePattern != `[.]go$` {
		t.Fatalf("host=%q filePattern=%v", host, filePattern)
	}
	if !reflect.DeepEqual(subcommandArgs, []string{"echo", "x"}) {
		t.Fatalf("subcommandArgs %#v", subcommandArgs)
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternMissingLongValue(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"--file-pattern"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "needs an argument: -file-pattern") {
		t.Fatalf("error %q", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternMissingShortValue(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"-f"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "needs an argument: -f") {
		t.Fatalf("error %q", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternEmptyEquals(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"--file-pattern=", "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "file-pattern value cannot be empty") {
		t.Fatalf("error %q", err.Error())
	}
}

func TestConsumeSubcommandSharedFlags_FilePatternEmptyNextToken(t *testing.T) {
	_, _, _, err := consumeSubcommandSharedFlags([]string{"-f", "", "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "file-pattern value cannot be empty") {
		t.Fatalf("error %q", err.Error())
	}
}

// TestInitializeContext_RunWithPatternGofmtTemplate covers initializeContext for argv after a shell parses:
//
//	run -f '[.]go$' gofmt -w __FILE_PATH__
//
// (Quotes are removed by the shell; each token is one os.Args element.) The regexp is stored on
// WatExecutionContext; remaining args are only the subprocess template (gofmt and its arguments).
func TestInitializeContext_RunWithPatternGofmtTemplate(t *testing.T) {
	wantPattern := `[.]go$`
	wantTemplate := []string{"gofmt", "-w", "__FILE_PATH__"}
	args := []string{"run", "-f", wantPattern, wantTemplate[0], wantTemplate[1], wantTemplate[2]}

	execCtx, subArgs, err := initializeContext(args)
	if err != nil {
		t.Fatalf("initializeContext: %v", err)
	}
	if execCtx.Subcommand() != "run" {
		t.Fatalf("Subcommand: want run, got %q", execCtx.Subcommand())
	}
	if execCtx.FilePattern() == nil || *execCtx.FilePattern() != wantPattern {
		t.Fatalf("FilePattern: want %q, got %#v", wantPattern, execCtx.FilePattern())
	}
	if !slices.Equal(subArgs, wantTemplate) {
		t.Fatalf("subcommand (subprocess template) argv: want %#v, got %#v", wantTemplate, subArgs)
	}
}
