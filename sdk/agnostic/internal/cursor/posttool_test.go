package cursor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func TestMapAfterFileEdit(t *testing.T) {
	ev := mapAfterFileEdit(sdkcursor.AfterFileEdit{
		FilePath: "/w/main.go",
		Edits: []sdkcursor.Edit{
			{OldString: "foo", NewString: "bar"},
		},
	})
	if ev.Tool == nil || ev.Tool.Name != tools.ToolEdit || ev.Tool.Native != sdkcursor.EventAfterFileEdit {
		t.Fatalf("tool = %#v", ev.Tool)
	}
	var input map[string]any
	if err := json.Unmarshal(ev.Tool.Input.Raw(), &input); err != nil {
		t.Fatal(err)
	}
	if input["file_path"] != "/w/main.go" {
		t.Fatalf("file_path = %v", input["file_path"])
	}
	if ev.Result == nil || string(ev.Result.Raw) == "" {
		t.Fatalf("result = %#v", ev.Result)
	}
}

func TestCallObservePostTool_DiscardsBuilderOutput(t *testing.T) {
	var saw model.PostToolEvent
	err := callObservePostTool(context.Background(), &model.PostToolEvent{
		Tool: &model.ToolCall{Name: tools.ToolEdit, Native: sdkcursor.EventAfterFileEdit},
	}, func(_ context.Context, event model.PostToolEvent, results model.PostToolResults) (model.PostToolResult, error) {
		saw = event
		out := results.Context("should not reach Cursor")
		if !out.IsZero() {
			t.Fatal("observe-only Context result must report IsZero")
		}
		out = out.WithUpdatedOutput("also ignored")
		if !out.IsZero() {
			t.Fatal("observe-only WithUpdatedOutput result must report IsZero")
		}
		return out, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if saw.Tool == nil || saw.Tool.Native != sdkcursor.EventAfterFileEdit {
		t.Fatalf("saw = %#v", saw)
	}
}

func TestRegisterPostTool_NilHandler(t *testing.T) {
	if RegisterPostTool(nil) != nil {
		t.Fatal("nil handler should return nil hooks")
	}
}
