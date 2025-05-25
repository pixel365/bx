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
		{"equals", fields{FooterTemplate: "some footer"}, "<br>some footer", false},
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

func TestChangelog_ApplyTransformation(t *testing.T) {
	type fields struct {
		Changelog Changelog
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		fields fields
	}{
		{"nil transformation", args{s: "raw string"}, "raw string", fields{Changelog: Changelog{}}},
		{
			"unknown type",
			args{s: " raw string "},
			"raw string",
			fields{Changelog: Changelog{Transform: &[]TypeValue[TransformType, []string]{
				{TransformType("unknown"), []string{""}},
			}}},
		},
		{
			"has prefix",
			args{s: "feat: some feature"},
			"some feature",
			fields{Changelog: Changelog{Transform: &[]TypeValue[TransformType, []string]{
				{StripPrefix, []string{"feat:"}},
			}}},
		},
		{
			"no prefix",
			args{s: "fix: some feature"},
			"fix: some feature",
			fields{Changelog: Changelog{Transform: &[]TypeValue[TransformType, []string]{
				{StripPrefix, []string{"feat:"}},
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Changelog.ApplyTransformation(tt.args.s); got != tt.want {
				t.Errorf("ApplyTransformation() = %v, want %v", got, tt.want)
			}
		})
	}
}
