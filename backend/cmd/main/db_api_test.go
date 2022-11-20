package main

import (
	"testing"
)

func TestDocument_String(t *testing.T) {
	tests := []struct {
		name    string
		obj     Document
		wantRes string
	}{
		{"empty", Document{}, ""},
		{"partial", Document{
			Name:       "имя",
			Bplace:     "место рождения",
			Years:      "годы",
			Vdate:      "2011-02-03",
			Vplace:     "место",
			Rang:       "звание",
			Awards:     "награды",
			PhDate:     "2015-05-04",
			Info:       "информация",
			SenderName: "заявитель",
			Phone:      "+7569874569",
		}, "[имя]\r\nимя" +
			"\r\n\r\n[место рождения]\r\nместо рождения" +
			"\r\n\r\n[годы жизни]\r\nгоды" +
			"\r\n\r\n[дата призыва]\r\n2011-02-03" +
			"\r\n\r\n[пункт призыва]\r\nместо" +
			"\r\n\r\n[воинское звание]\r\nзвание" +
			"\r\n\r\n[награды]\r\nнаграды" +
			"\r\n\r\n[когда сделано фото]\r\n2015-05-04" +
			"\r\n\r\n[доп. сведения]\r\nинформация" +
			"\r\n\r\n[контакт фио]\r\nзаявитель" +
			"\r\n\r\n[контакт телефон]\r\n+7569874569" +
			"\r\n\r\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.obj.String(); gotRes != tt.wantRes {
				t.Errorf("String() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
