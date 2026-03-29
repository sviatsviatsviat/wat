package cursor

import "testing"

func TestNewHookDataAfterFileEdit_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "conv-1",
		"file_path": "C:/repo/file.txt",
		"edits": [
			{"old_string":"old-1","new_string":"new-1"},
			{"old_string":"old-2","new_string":"new-2"}
		]
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterFileEdit", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "C:/repo/file.txt", hookData.Fields.FilePath)
	if len(hookData.Fields.Edits) != 2 {
		t.Fatalf("Edits length: want 2, got %d", len(hookData.Fields.Edits))
	}
	assertEqual(t, "old-1", hookData.Fields.Edits[0].OldString)
	assertEqual(t, "new-1", hookData.Fields.Edits[0].NewString)
	assertEqual(t, "old-2", hookData.Fields.Edits[1].OldString)
	assertEqual(t, "new-2", hookData.Fields.Edits[1].NewString)
}

func TestNewHookDataAfterFileEdit_zeroValueFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"file_path": "",
		"edits": []
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	assertEqual(t, "", hookData.Fields.FilePath)
	assertIntEqual(t, 0, len(hookData.Fields.Edits))
}

func TestNewHookDataAfterFileEdit_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewHookDataAfterTabFileEdit_fullPayloadWithRangeAndLines(t *testing.T) {
	input := `{
		"hook_event_name": "afterTabFileEdit",
		"conversation_id": "conv-tab",
		"file_path": "C:/repo/file.go",
		"edits": [
			{
				"old_string": "search",
				"new_string": "replace",
				"range": {
					"start_line_number": 10,
					"start_column": 5,
					"end_line_number": 10,
					"end_column": 20
				},
				"old_line": "line before edit",
				"new_line": "line after edit"
			}
		]
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterTabFileEdit", hookData.HookEventName)
	assertEqual(t, "conv-tab", hookData.ConversationID)
	assertEqual(t, "C:/repo/file.go", hookData.Fields.FilePath)
	if len(hookData.Fields.Edits) != 1 {
		t.Fatalf("Edits length: want 1, got %d", len(hookData.Fields.Edits))
	}
	e := hookData.Fields.Edits[0]
	assertEqual(t, "search", e.OldString)
	assertEqual(t, "replace", e.NewString)
	assertEqual(t, "line before edit", e.OldLine)
	assertEqual(t, "line after edit", e.NewLine)
	if e.EditRange == nil {
		t.Fatal("EditRange: want non-nil")
	}
	assertIntEqual(t, 10, e.EditRange.StartLineNumber)
	assertIntEqual(t, 5, e.EditRange.StartColumn)
	assertIntEqual(t, 10, e.EditRange.EndLineNumber)
	assertIntEqual(t, 20, e.EditRange.EndColumn)
}
