package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UploadResult(t *testing.T) {
	t.Parallel()

	type args struct {
		htmlContent string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"has error", args{htmlContent: `<p class="paragraph-15 color-red m-0">error</p>`}, true},
		{"no error", args{htmlContent: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := UploadResult(tt.args.htmlContent)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseSessionId(t *testing.T) {
	t.Parallel()

	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"success", args{`<input type="hidden" name="sessid" id="sessid" value="xxx" />`}, "xxx"},
		{"not found", args{`<p>some paragraph</p>`}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ParseSessionId(tt.args.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseVersions_success(t *testing.T) {
	t.Parallel()

	content := `<table class="data-table mt-3 mb-3">
<tbody>
<tr>
<td>1.0.1</td>
<td>
<input type="radio" name="1.0.1" id="1.0.1alpha" value="alpha" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1alpha" title="Alpha">Alpha</label><br>
<input type="radio" name="1.0.1" id="1.0.1beta" value="beta" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1beta" title="Beta">Beta</label><br>
<input type="radio" name="1.0.1" id="1.0.1stable" value="stable" checked="" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1stable" title="Доступно всем клиентам">Stable</label><br>
</td>
</tr>
</tbody>
</table>`

	versions, err := ParseVersions(content)
	require.NoError(t, err)
	assert.Len(t, versions, 1)
}

func TestParseVersions_table_not_found(t *testing.T) {
	t.Parallel()

	content := `<table>
<tbody>
<tr>
<td>1.0.1</td>
<td>
<input type="radio" name="1.0.1" id="1.0.1alpha" value="alpha" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1alpha" title="Alpha">Alpha</label><br>
<input type="radio" name="1.0.1" id="1.0.1beta" value="beta" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1beta" title="Beta">Beta</label><br>
<input type="radio" name="1.0.1" id="1.0.1stable" value="stable" checked="" onclick="changeType('1.0.1', this.value)">
<label for="1.0.1stable" title="Доступно всем клиентам">Stable</label><br>
</td>
</tr>
</tbody>
</table>`

	_, err := ParseVersions(content)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "table not found")
}

func TestParseVersions_invalid_content(t *testing.T) {
	t.Parallel()

	_, err := ParseVersions(`<table>`)
	require.Error(t, err)
}
