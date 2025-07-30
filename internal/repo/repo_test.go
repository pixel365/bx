package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types/changelog"

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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOpenRepository_Ok(t *testing.T) {
	pwd, _ := filepath.Abs("../../")
	_, err := OpenRepository(pwd)
	require.NoError(t, err)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
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
	require.NoError(t, err)
	assert.Equal(t, []string{"commit 1", "commit 2"}, commits)
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
	require.NoError(t, err)
	assert.Equal(t, []string{"commit 2", "commit 1"}, commits)
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
			got := CommitFilter(tt.args.message, tt.args.conditions)
			assert.Equal(t, tt.want, got)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChangesList_nil_repository(t *testing.T) {
	origOpenRepositoryFunc := openRepositoryFunc
	defer func() { openRepositoryFunc = origOpenRepositoryFunc }()

	openRepositoryFunc = func(_ string) (*git.Repository, error) {
		return nil, nil
	}

	_, err := ChangesList("", changelog.Changelog{})
	assert.ErrorIs(t, errors2.ErrNilRepository, err)
}

func TestChangesList_hashes_fail(t *testing.T) {
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
	require.Error(t, err)
	assert.Equal(t, e, err.Error())
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
			got := o.IsChangedFile(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
