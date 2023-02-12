package applemusic

import (
	"testing"
	"time"
)

func Test_releaseYearWithinRange(t *testing.T) {
	type args struct {
		releaseYearCandidate int
		releaseYearInput     int
		rangeYears           int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test current year",
			args: args{
				releaseYearCandidate: 2023,
				releaseYearInput:     2023,
				rangeYears:           3,
			},
			want: true,
		},
		{
			name: "test outside of range",
			args: args{
				releaseYearCandidate: time.Now().Year() + 4,
				releaseYearInput:     2023,
				rangeYears:           3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := releaseYearWithinRange(tt.args.releaseYearCandidate, tt.args.releaseYearInput, tt.args.rangeYears); got != tt.want {
				t.Errorf("releaseYearWithinRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
