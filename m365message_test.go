package m365message_test

import (
	"testing"

	m365message "github.com/alexhowarth/go-m365-msg-builder"
)

func TestMessage_Build(t *testing.T) {
	type fields struct {
		rw        m365message.ReadWrite
		payload   []int
		position  int
		direction m365message.Direction
	}
	tests := []struct {
		name      string
		fields    fields
		want      string
		wantErr   bool
		wantedErr error
	}{
		// handle missing args
		{"ReadWrite", fields{0, nil, 0, 0}, "", true, m365message.ErrorInvalidReadWrite},
		{"Direction", fields{m365message.READ, nil, 0, 0}, "", true, m365message.ErrorInvalidDirection},
		{"Payload", fields{m365message.READ, nil, 0, m365message.MASTER_TO_M365}, "", true, m365message.ErrorPayloadRequired},
		{"Position", fields{m365message.READ, []int{0x02}, 0, m365message.MASTER_TO_M365}, "", true, m365message.ErrorPositionRequired},
		// test output
		{"Build", fields{m365message.READ, []int{0x02}, 0x1A, m365message.MASTER_TO_M365}, "55aa0320011a02bfff", false, nil},
		{"Light On", fields{m365message.WRITE, []int{0x0002}, 0x7D, m365message.MASTER_TO_M365}, "55aa0320037d025aff", false, nil},
		{"Battery Stats", fields{m365message.READ, []int{0x0a}, 0x31, m365message.MASTER_TO_BATTERY}, "55aa032201310a9eff", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := m365message.NewMessage()
			m.SetRW(tt.fields.rw)
			m.SetPayload(tt.fields.payload)
			m.SetPosition(tt.fields.position)
			m.SetDirection(tt.fields.direction)
			got, err := m.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("m365message.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && tt.wantedErr != err {
				t.Errorf("m365message.Build() error = %v, wantErr %v, wantedErr %v", err, tt.wantErr, tt.wantedErr)
				return
			}
			if got != tt.want {
				t.Errorf("m365message.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
