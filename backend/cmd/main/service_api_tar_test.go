package main

import (
	"testing"
	"time"
)

func Test_parseDate(t *testing.T) {
	zeroTime := time.Unix(0, 0)
	type args struct {
		filename string
		location string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"empty",
			args{"", ""},
			zeroTime,
			true},
		{"bad location",
			args{"filename.ext", "Eur0pe/M"},
			zeroTime,
			true},
		{"bad filename",
			args{"filename.ext", "Europe/Moscow"},
			zeroTime,
			true},
		{"ok",
			args{"2012-05-08.ext", "Europe/Moscow"},
			time.Unix(1336410000, 0).UTC().Add(time.Hour * 4),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.filename, tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Equal(tt.want) {
				t.Errorf("parseDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
