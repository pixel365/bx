package internal

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5"
)

func TestOpenRepository(t *testing.T) {
	type args struct {
		repository string
	}
	tests := []struct {
		want    *git.Repository
		name    string
		args    args
		wantErr bool
	}{
		{nil, "empty repository", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenRepository(tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenRepository_Ok(t *testing.T) {
	t.Run("repository exists", func(t *testing.T) {
		pwd, _ := filepath.Abs("../")
		_, err := OpenRepository(pwd)
		if err != nil {
			t.Errorf("OpenRepository() error = %v", err)
		}
	})
}

func TestChangelogList(t *testing.T) {
	type args struct {
		repository string
		rules      Changelog
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"empty changelog", args{"", Changelog{}}, []string{}, false},
		{"empty repository", args{"", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
			Sort:      "asc",
			Condition: TypeValue[ChangelogConditionType, []string]{},
		}}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangelogList(tt.args.repository, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangelogList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangelogList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommitFilter(t *testing.T) {
	type args struct {
		message    string
		conditions TypeValue[ChangelogConditionType, []string]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"allow", args{
			message: "feat: some feature",
			conditions: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat:([\W\w]+)$`,
				},
			},
		}, true},
		{"not allow", args{
			message: "fix: some fix",
			conditions: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat:([\W\w]+)$`,
				},
			},
		}, false},
		{"exclude", args{
			message: "fix: some fix",
			conditions: TypeValue[ChangelogConditionType, []string]{
				Type: Exclude,
				Value: []string{
					`^feat:([\W\w]+)$`,
				},
			},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommitFilter(tt.args.message, tt.args.conditions); got != tt.want {
				t.Errorf("CommitFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listOfCommits(t *testing.T) {
	type args struct {
		repository *git.Repository
		filter     CommitFilterFunc
		rules      Changelog
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"nil repository", args{}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listOfCommits(tt.args.repository, tt.args.rules, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("listOfCommits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listOfCommits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashes(t *testing.T) {
	type args struct {
		repository *git.Repository
		rules      Changelog
	}
	tests := []struct {
		name    string
		args    args
		want    plumbing.Hash
		want1   plumbing.Hash
		wantErr bool
	}{
		{"nil repository", args{}, plumbing.ZeroHash, plumbing.ZeroHash, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := hashes(tt.args.repository, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hashes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("hashes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChangesList(t *testing.T) {
	type args struct {
		repository string
		rules      Changelog
	}
	tests := []struct {
		want    *Changes
		name    string
		args    args
		wantErr bool
	}{
		{nil, "empty repository", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangesList(tt.args.repository, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangesList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChanges_IsChangedFile(t *testing.T) {
	type fields struct {
		Added    []string
		Modified []string
		Deleted  []string
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
		{"added", args{path: "test.txt"}, fields{Added: []string{"test.txt"}}, true},
		{"modified", args{path: "modified.txt"}, fields{Modified: []string{"modified.txt"}}, true},
		{"deleted", args{path: "modified.txt"}, fields{Deleted: []string{"test.txt"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Changes{
				Added:    tt.fields.Added,
				Modified: tt.fields.Modified,
				Deleted:  tt.fields.Deleted,
			}
			if got := o.IsChangedFile(tt.args.path); got != tt.want {
				t.Errorf("IsChangedFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
