package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pixel365/bx/internal/types/changelog"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"

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
		pwd, _ := filepath.Abs("../../")
		_, err := OpenRepository(pwd)
		if err != nil {
			t.Errorf("OpenRepository() error = %v", err)
		}
	})
}

func TestChangelogList(t *testing.T) {
	type args struct {
		repository string
		rules      changelog.Changelog
	}
	tests := []struct {
		name    string
		want    []string
		args    args
		wantErr bool
	}{
		{"empty changelog", []string{}, args{"", changelog.Changelog{}}, false},
		{"empty repository", nil, args{"", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Sort:      "asc",
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{},
		}}, true},
		{"empty from values", []string{}, args{"", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
			Sort:      "asc",
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{},
		}}, false},
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

func TestChangelogList_listOfCommits_Fail(t *testing.T) {
	origListOfCommits := listOfCommitsFunc
	origOpenRepository := openRepositoryFunc
	defer func() {
		listOfCommitsFunc = origListOfCommits
	}()
	defer func() {
		openRepositoryFunc = origOpenRepository
	}()

	listOfCommitsFunc = func(_ *git.Repository, _ changelog.Changelog, _ CommitFilterFunc) ([]string, error) {
		return nil, errors.New("fail")
	}

	openRepositoryFunc = func(_ string) (*git.Repository, error) {
		return &git.Repository{}, nil
	}

	t.Run("commits fail", func(t *testing.T) {
		_, err := ChangelogList("", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		})
		if err == nil {
			t.Errorf("ChangelogList() error = %v, wantErr %v", err, errors.New("fail"))
		}
	})
}

func TestChangelogList_listOfCommits_Ok_Asc(t *testing.T) {
	origListOfCommits := listOfCommitsFunc
	origOpenRepository := openRepositoryFunc
	defer func() {
		listOfCommitsFunc = origListOfCommits
	}()
	defer func() {
		openRepositoryFunc = origOpenRepository
	}()

	listOfCommitsFunc = func(_ *git.Repository, _ changelog.Changelog, _ CommitFilterFunc) ([]string, error) {
		return []string{"commit 2", "commit 1"}, nil
	}

	openRepositoryFunc = func(_ string) (*git.Repository, error) {
		return &git.Repository{}, nil
	}

	t.Run("commits fail", func(t *testing.T) {
		commits, err := ChangelogList("", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Sort: types.Asc,
		})
		if err != nil {
			t.Errorf("ChangelogList() error = %v", err)
		}

		if commits[0] != "commit 1" {
			t.Errorf("ChangelogList() got = %v, want %v", commits[0], "commit 1")
		}
	})
}

func TestChangelogList_listOfCommits_Ok_Desc(t *testing.T) {
	origListOfCommits := listOfCommitsFunc
	origOpenRepository := openRepositoryFunc
	defer func() {
		listOfCommitsFunc = origListOfCommits
	}()
	defer func() {
		openRepositoryFunc = origOpenRepository
	}()

	listOfCommitsFunc = func(_ *git.Repository, _ changelog.Changelog, _ CommitFilterFunc) ([]string, error) {
		return []string{"commit 1", "commit 2"}, nil
	}

	openRepositoryFunc = func(_ string) (*git.Repository, error) {
		return &git.Repository{}, nil
	}

	t.Run("commits fail", func(t *testing.T) {
		commits, err := ChangelogList("", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Sort: types.Desc,
		})
		if err != nil {
			t.Errorf("ChangelogList() error = %v", err)
		}

		if commits[0] != "commit 2" {
			t.Errorf("ChangelogList() got = %v, want %v", commits[0], "commit 2")
		}
	})
}

func TestCommitFilter(t *testing.T) {
	type args struct {
		message    string
		conditions types.TypeValue[types.ChangelogConditionType, []string]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"allow", args{
			message: "feat: some feature",
			conditions: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat:([\W\w]+)$`,
				},
			},
		}, true},
		{"not allow", args{
			message: "fix: some fix",
			conditions: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat:([\W\w]+)$`,
				},
			},
		}, false},
		{"exclude", args{
			message: "fix: some fix",
			conditions: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Exclude,
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
		rules      changelog.Changelog
	}
	tests := []struct {
		name    string
		want    []string
		args    args
		wantErr bool
	}{
		{"nil repository", nil, args{}, true},
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
		rules      changelog.Changelog
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
		rules      changelog.Changelog
	}
	tests := []struct {
		want    *types.Changes
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

func TestChangesList_nil_repository(t *testing.T) {
	t.Run("nil repository", func(t *testing.T) {
		origOpenRepositoryFunc := openRepositoryFunc
		defer func() { openRepositoryFunc = origOpenRepositoryFunc }()

		openRepositoryFunc = func(_ string) (*git.Repository, error) {
			return nil, nil
		}

		_, err := ChangesList("", changelog.Changelog{})
		if !errors.Is(err, errors2.ErrNilRepository) {
			t.Errorf("ChangesList() error = %v, wantErr %v", err, errors2.ErrNilRepository)
		}
	})
}

func TestChangesList_hashes_fail(t *testing.T) {
	t.Run("hashes fail", func(t *testing.T) {
		origOpenRepositoryFunc := openRepositoryFunc
		origHashesFunc := hashesFunc
		defer func() { openRepositoryFunc = origOpenRepositoryFunc }()
		defer func() { hashesFunc = origHashesFunc }()

		openRepositoryFunc = func(_ string) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		hashesFunc = func(_ *git.Repository, _ changelog.Changelog) (plumbing.Hash, plumbing.Hash, error) {
			return plumbing.ZeroHash, plumbing.ZeroHash, errors.New("some error")
		}

		_, err := ChangesList("repo", changelog.Changelog{})
		e := fmt.Errorf("repository [%s]: %w", "repo", errors.New("some error")).Error()
		if err.Error() != e {
			t.Errorf("ChangesList() error = %v, want %v", e, err)
		}
	})
}

func TestChanges_IsChangedFile(t *testing.T) {
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
		{"added", args{path: "test.txt"}, fields{Added: []string{"test.txt"}}, true},
		{"modified", args{path: "modified.txt"}, fields{Modified: []string{"modified.txt"}}, true},
		{"deleted", args{path: "modified.txt"}, fields{Deleted: []string{"test.txt"}}, false},
		{"moved", args{path: "modified.txt"}, fields{Moved: []string{"modified.txt"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &types.Changes{
				Added:    tt.fields.Added,
				Modified: tt.fields.Modified,
				Deleted:  tt.fields.Deleted,
				Moved:    tt.fields.Moved,
			}
			if got := o.IsChangedFile(tt.args.path); got != tt.want {
				t.Errorf("IsChangedFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
