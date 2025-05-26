package changelog

import (
	"fmt"
	"testing"

	"github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"
)

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
			fields{
				Changelog: Changelog{Transform: &[]types.TypeValue[types.TransformType, []string]{
					{Type: types.TransformType("unknown"), Value: []string{""}},
				}},
			},
		},
		{
			"has prefix",
			args{s: "feat: some feature"},
			"some feature",
			fields{
				Changelog: Changelog{Transform: &[]types.TypeValue[types.TransformType, []string]{
					{Type: types.StripPrefix, Value: []string{"feat:"}},
				}},
			},
		},
		{
			"no prefix",
			args{s: "fix: some feature"},
			"fix: some feature",
			fields{
				Changelog: Changelog{Transform: &[]types.TypeValue[types.TransformType, []string]{
					{Type: types.StripPrefix, Value: []string{"feat:"}},
				}},
			},
		},
		{
			"no suffix",
			args{s: "some feature fix"},
			"some feature fix",
			fields{
				Changelog: Changelog{Transform: &[]types.TypeValue[types.TransformType, []string]{
					{Type: types.StripSuffix, Value: []string{"feat"}},
				}},
			},
		},
		{
			"has suffix",
			args{s: "some feature fix"},
			"some feature",
			fields{
				Changelog: Changelog{Transform: &[]types.TypeValue[types.TransformType, []string]{
					{Type: types.StripSuffix, Value: []string{"fix"}},
				}},
			},
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

func Test_transformValidate(t *testing.T) {
	type args struct {
		transform *[]types.TypeValue[types.TransformType, []string]
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{transform: nil}, "nil transform", false},
		{
			args{transform: &[]types.TypeValue[types.TransformType, []string]{}},
			"empty transform",
			false,
		},
		{args{transform: &[]types.TypeValue[types.TransformType, []string]{
			{Type: types.TransformType("unknown"), Value: []string{"value"}},
		}}, "unknown transform", true},
		{args{transform: &[]types.TypeValue[types.TransformType, []string]{
			{Type: types.StripPrefix, Value: []string{""}},
		}}, "empty value", true},
		{args{transform: &[]types.TypeValue[types.TransformType, []string]{
			{Type: types.StripPrefix, Value: []string{"feat:"}},
		}}, "valid transform", false},
		{args{transform: &[]types.TypeValue[types.TransformType, []string]{
			{Type: types.StripPrefix, Value: []string{}},
		}}, "empty values", true},
		{args{transform: &[]types.TypeValue[types.TransformType, []string]{
			{Type: types.StripSuffix, Value: []string{"feat:"}},
		}}, "valid transform", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := transformValidate(tt.args.transform); (err != nil) != tt.wantErr {
				t.Errorf("transformValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_changeLogFromToValidate(t *testing.T) {
	type args struct {
		c *Changelog
	}
	tests := []struct {
		wantErr error
		args    args
		name    string
	}{
		{
			name: "valid commit to tag",
			args: args{
				c: &Changelog{
					From: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Commit,
						Value: "abc123",
					},
					To: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Tag,
						Value: "v1.2.3",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "missing from value",
			args: args{
				c: &Changelog{
					From: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Commit,
						Value: "",
					},
					To: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Tag,
						Value: "v1.2.3",
					},
				},
			},
			wantErr: errors.ErrChangelogValue,
		},
		{
			name: "missing to value",
			args: args{
				c: &Changelog{
					From: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Commit,
						Value: "abc123",
					},
					To: types.TypeValue[types.ChangelogType, string]{Type: types.Tag, Value: ""},
				},
			},
			wantErr: errors.ErrChangelogValue,
		},
		{
			name: "invalid from type",
			args: args{
				c: &Changelog{
					From: types.TypeValue[types.ChangelogType, string]{
						Type:  "branch",
						Value: "abc123",
					},
					To: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Tag,
						Value: "v1.2.3",
					},
				},
			},
			wantErr: fmt.Errorf("changelog from: type must be %s or %s", types.Commit, types.Tag),
		},
		{
			name: "invalid to type",
			args: args{
				c: &Changelog{
					From: types.TypeValue[types.ChangelogType, string]{
						Type:  types.Commit,
						Value: "abc123",
					},
					To: types.TypeValue[types.ChangelogType, string]{
						Type:  "branch",
						Value: "v1.2.3",
					},
				},
			},
			wantErr: fmt.Errorf("changelog to: type must be %s or %s", types.Commit, types.Tag),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := changeLogFromToValidate(tt.args.c)
			if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("unexpected error.\n want: %q\n  got: %q", tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestChangelog_IsValid(t *testing.T) {
	type tv = types.TypeValue[types.ChangelogType, string]
	type fields struct {
		Changelog *Changelog
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{Changelog: &Changelog{}}, "empty changelog", true},
		{fields{Changelog: &Changelog{
			From: tv{Type: types.Tag, Value: "tag1"},
			To:   tv{Type: types.Tag, Value: "tag2"},
			Sort: "",
		},
		}, "empty sort", false},
		{fields{Changelog: &Changelog{
			From: tv{Type: types.Tag, Value: "tag1"},
			To:   tv{Type: types.Tag, Value: "tag2"},
			Sort: types.Desc,
		},
		}, "valid sort", false},
		{fields{Changelog: &Changelog{
			From: tv{Type: types.Tag, Value: "tag1"},
			To:   tv{Type: types.Tag, Value: "tag2"},
			Sort: types.SortingType("unknown"),
		},
		}, "invalid sort", true},
		{fields{Changelog: &Changelog{
			From: tv{Type: types.Tag, Value: "tag1"},
			To:   tv{Type: types.Tag, Value: "tag2"},
			Sort: types.Desc,
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
			},
		},
		}, "invalid condition", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.Changelog.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_conditionValidate(t *testing.T) {
	type tv types.TypeValue[types.ChangelogConditionType, []string]
	tests := []struct {
		name      string
		condition tv
		wantErr   bool
	}{
		{"empty condition", tv{Type: ""}, false},
		{"invalid condition", tv{Type: types.ChangelogConditionType("unknown")}, true},
		{"empty values", tv{Type: types.Include}, true},
		{"empty value", tv{Type: types.Include, Value: []string{""}}, true},
		{"invalid regex", tv{Type: types.Include, Value: []string{"[unclosed"}}, true},
		{"valid regex", tv{Type: types.Include, Value: []string{"^feat:([\\W\\w]+)$"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := types.TypeValue[types.ChangelogConditionType, []string](tt.condition)
			if err := conditionValidate(cond); (err != nil) != tt.wantErr {
				t.Errorf("conditionValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
