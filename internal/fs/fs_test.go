package fs

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/types"
)

type FakeModuleConfig struct{}

func (f FakeModuleConfig) GetVariables() map[string]string { return nil }
func (f FakeModuleConfig) GetRun() map[string][]string     { return nil }
func (f FakeModuleConfig) GetStages() []types.Stage        { return nil }
func (f FakeModuleConfig) GetIgnore() []string {
	return []string{
		"**/*.log",
		"*.json",
		"**/*some*/*",
	}
}
func (f FakeModuleConfig) GetChanges() *types.Changes { return nil }
func (f FakeModuleConfig) IsLastVersion() bool        { return false }

type FakeFileInfo struct {
	Dir bool
}

func (f FakeFileInfo) Name() string       { return "" }
func (f FakeFileInfo) Size() int64        { return 0 }
func (f FakeFileInfo) Mode() fs.FileMode  { return fs.ModeDir }
func (f FakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f FakeFileInfo) Sys() any           { return nil }
func (f FakeFileInfo) IsDir() bool        { return f.Dir }

func Test_mkdir(t *testing.T) {
	name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
	path, err := MkDir(name)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(path)
	})
}

func Test_zipIt(t *testing.T) {
	name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
	path, err := MkDir(name)
	require.NoError(t, err)

	archivePath := fmt.Sprintf("./_%d.zip", time.Now().UTC().Unix())
	err = ZipIt(path, archivePath)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(path)
		_ = os.Remove(archivePath)
	})
}

func Test_shouldSkip(t *testing.T) {
	t.Parallel()
	patterns := []string{
		"**/*.log",
		"*.json",
		"**/*some*/*",
	}
	type args struct {
		path     string
		patterns []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{".", nil}, false},
		{"2", args{"./testing/errors.log", patterns}, true},
		{"3", args{"./testing/errors.json", patterns}, false},
		{"4", args{"errors.json", patterns}, true},
		{"5", args{"./testing/data/awesome/cfg.yaml", patterns}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldSkip(tt.args.path, tt.args.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_CopyFromPath_ok(t *testing.T) {
	from := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
	fromPath, err := MkDir(from)
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(fromPath); err != nil {
			t.Error(err)
		}
	}()

	to := fmt.Sprintf("./__%d", time.Now().UTC().Unix())
	toPath, err := MkDir(to)
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(toPath); err != nil {
			t.Error(err)
		}
	}()

	fileName := fmt.Sprintf("%d.txt", time.Now().UTC().Unix())
	filePath := filepath.Join(from, fileName)
	filePath = filepath.Clean(filePath)
	file, err := os.Create(filePath)
	require.NoError(t, err)

	err = file.Close()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Error(err)
		}
	}()

	errChan := make(chan types.Path, 1)

	module := FakeModuleConfig{}

	path := types.Path{
		From:           from,
		To:             to,
		ActionIfExists: types.Replace,
		Convert:        false,
	}

	if err = PathProcessing(
		context.Background(),
		errChan,
		&module,
		path,
		[]string{},
	); err != nil {
		close(errChan)
		t.Error(err)
	}

	defer func() {
		_ = os.Remove(fmt.Sprintf("%s/%s", toPath, fileName))
	}()

	close(errChan)
}

func TestPathProcessingContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := PathProcessing(ctx, nil, nil, types.Path{}, nil)
	require.Error(t, err)
}

func Test_isConvertable(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{args{"/some/lang/file.php"}, "php", true},
		{args{"/some/path/file.php"}, "php", false},
		{args{"/some/path/description.ru"}, "description.ru", true},
		{args{"/some/path/image.jpg"}, "jpg", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isConvertable(tt.args.path)
			if tt.want {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}
		})
	}
}

func Test_isEmptyDir(t *testing.T) {
	name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
	path, err := MkDir(name)
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(path); err != nil {
			t.Error(err)
		}
	}()

	assert.True(t, IsEmptyDir(path))
}

func Test_removeEmptyDirs(t *testing.T) {
	name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
	name2 := fmt.Sprintf("./%s/%d", name, time.Now().UTC().Unix())
	path, err := MkDir(name)
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(path); err != nil {
			t.Error(err)
		}
	}()

	path2, err := MkDir(name2)
	require.NoError(t, err)

	status, err := RemoveEmptyDirs(path)
	require.NoError(t, err)
	assert.True(t, status)

	if !status || err != nil {
		defer func() {
			if err := os.Remove(path2); err != nil {
				t.Error(err)
			}
		}()
	}
}

func Test_shouldInclude(t *testing.T) {
	t.Parallel()
	type args struct {
		path     string
		patterns []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty patterns", args{"./testing.php", []string{}}, true},
		{"empty path", args{"", []string{"**/*.php"}}, true},
		{"included path", args{"./testing.php", []string{"**/*.php"}}, true},
		{"excluded json", args{"./testing.json", []string{"!**/*.json"}}, false},
		{"included php", args{"./testing.php", []string{"!**/*.json"}}, true},
		{
			"excluded test file",
			args{"./some_test.php", []string{"**/*.php", "!**/*_test.php"}},
			false,
		},
		{
			"mutually exclusive rules",
			args{"./testing.php", []string{"**/*.php", "!**/*.php"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldInclude(tt.args.path, tt.args.patterns)
			if tt.want {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}
		})
	}
}

func TestIsFileExists_true(t *testing.T) {
	filePath := fmt.Sprintf("./%d.txt", time.Now().UTC().Unix())
	filePath = filepath.Clean(filePath)
	file, err := os.Create(filePath)
	require.NoError(t, err)

	_, err = file.WriteString("str")
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Error(err)
		}
	}()

	ok, size := IsFileExists(filePath)
	assert.False(t, !ok || size == 0)
}

func TestIsFileExists_false(t *testing.T) {
	ok, size := IsFileExists("./some-file.txt")
	assert.False(t, ok || size > 0)
}

func Test_skip(t *testing.T) {
	t.Parallel()
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{info: nil}, "nil info", false},
		{args{info: FakeFileInfo{Dir: true}}, "is dir", true},
		{args{info: FakeFileInfo{Dir: false}}, "is not dir", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := skip(tt.args.info)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_visitor(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		cfg         interfaces.ModuleConfig
		err         error
		filesCh     chan<- types.Path
		changes     *types.Changes
		path        types.Path
		filterRules []string
	}
	tests := []struct {
		want filepath.WalkFunc
		name string
		args args
	}{
		{want: func(_ string, _ fs.FileInfo, _ error) error {
			return errors.ErrNilContext
		}, name: "nil context", args: args{cfg: FakeModuleConfig{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			visit := visitor(
				ctx,
				tt.args.filesCh,
				tt.args.cfg,
				tt.args.path,
				tt.args.filterRules,
				tt.args.changes,
			)
			assert.NotNil(t, visit)

			err := visit("", FakeFileInfo{}, tt.args.err)
			if err != nil {
				assert.NotErrorIs(t, context.Canceled, err)
			}
		})
	}
}
