package main

import (
	"testing"
)

func Test_scaleSizes(t *testing.T) {
	type args struct {
		w  uint
		h  uint
		to uint
	}
	tests := []struct {
		name     string
		args     args
		wantNewW uint
		wantNewH uint
	}{
		{"empty", args{0, 0, 0}, 0, 0},
		{"ok", args{100, 25, 50}, 50, 12},
		{"ok", args{25, 100, 50}, 12, 50},
		{"ok", args{25, 25, 50}, 25, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewW, gotNewH := scaleSizes(tt.args.w, tt.args.h, tt.args.to)
			if gotNewW != tt.wantNewW {
				t.Errorf("scaleSizes() gotNewW = %v, want %v", gotNewW, tt.wantNewW)
			}
			if gotNewH != tt.wantNewH {
				t.Errorf("scaleSizes() gotNewH = %v, want %v", gotNewH, tt.wantNewH)
			}
		})
	}
}
