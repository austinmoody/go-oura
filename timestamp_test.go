package go_oura

import (
	"testing"
	"time"
)

func TestTimeStamp_toTime(t *testing.T) {
	tests := []struct {
		name            string
		timeStampString TimeStamp
		want            time.Time
		wantErr         bool
	}{
		{
			name:            "valid ISO 8601 timestamp",
			timeStampString: "2024-01-01T15:04:05+00:00",
			want:            time.Date(2024, time.January, 1, 15, 4, 5, 0, time.UTC),
			wantErr:         false,
		},
		{
			name:            "valid RFC3339 timestamp",
			timeStampString: "2024-01-06T16:43:56Z",
			want:            time.Date(2024, time.January, 6, 16, 43, 56, 0, time.UTC),
			wantErr:         false,
		},
		{
			name:            "empty timestamp",
			timeStampString: "",
			want:            time.Time{},
			wantErr:         true,
		},
		{
			name:            "invalid timestamp",
			timeStampString: "invalid",
			want:            time.Time{},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.timeStampString.toTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("toTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("toTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
