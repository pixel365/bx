package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChanges_IsChangedFile(t *testing.T) {
	t.Parallel()

	type fields struct {
		Added    []string
		Modified []string
		Deleted  []string
		Moved    []string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   bool
	}{
		{
			args:   args{path: "some-added.txt"},
			fields: fields{Added: []string{"some-added.txt"}},
			want:   true,
		},
		{
			args:   args{path: "some-modified.txt"},
			fields: fields{Modified: []string{"some-modified.txt"}},
			want:   true,
		},
		{
			args:   args{path: "some-deleted.txt"},
			fields: fields{Deleted: []string{"some-deleted.txt"}},
			want:   false,
		},
		{
			args:   args{path: "some-moved.txt"},
			fields: fields{Moved: []string{"some-moved.txt"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := &Changes{
				Added:    tt.fields.Added,
				Modified: tt.fields.Modified,
				Deleted:  tt.fields.Deleted,
				Moved:    tt.fields.Moved,
			}
			got := o.IsChangedFile(tt.args.path)
			if tt.want {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}
		})
	}
}
