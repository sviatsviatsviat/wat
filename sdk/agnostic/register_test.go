package agnostic

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

func TestHandlerErrorExit(t *testing.T) {
	tests := []struct {
		name string
		ev   *model.Event
		want int
	}{
		{
			name: "copilot pre tool",
			ev:   &model.Event{Agent: model.Copilot, Kind: model.KindPreTool},
			want: 1,
		},
		{
			name: "copilot stop",
			ev:   &model.Event{Agent: model.Copilot, Kind: model.KindStop},
			want: 1,
		},
		{
			name: "cursor",
			ev:   &model.Event{Agent: model.Cursor, Kind: model.KindPreTool},
			want: cursor.HandlerErrorExit,
		},
		{
			name: "claude",
			ev:   &model.Event{Agent: model.Claude, Kind: model.KindPreTool},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlerErrorExit(tt.ev); got != tt.want {
				t.Fatalf("handlerErrorExit() = %d, want %d", got, tt.want)
			}
		})
	}
}
