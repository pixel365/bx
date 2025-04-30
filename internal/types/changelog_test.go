package types

import "testing"

func TestChangelog_EncodedFooter(t *testing.T) {
	type fields struct {
		FooterTemplate string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"empty", fields{FooterTemplate: ""}, "", false},
		{"equals", fields{FooterTemplate: "some footer"}, "\nsome footer", false},
		{"es", fields{FooterTemplate: "algún pie de página"}, "", true},
		{"jp", fields{FooterTemplate: "フッター"}, "", true},
		{"kz", fields{FooterTemplate: "кейбір төменгі колонтитул"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Changelog{
				FooterTemplate: tt.fields.FooterTemplate,
			}
			got, err := c.EncodedFooter()
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodedFooter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodedFooter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
