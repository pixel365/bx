package parser

import (
	"testing"
)

func Test_UploadResult(t *testing.T) {
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
			if err := UploadResult(tt.args.htmlContent); (err != nil) != tt.wantErr {
				t.Errorf("UploadResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseSessionId(t *testing.T) {
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
			if got := ParseSessionId(tt.args.content); got != tt.want {
				t.Errorf("ParseSessionId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVersions_success(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		versions, err := ParseVersions(content)
		if err != nil {
			t.Error(err)
		}
		if len(versions) != 1 {
			t.Errorf("len(versions) = %d, want 1", len(versions))
		}
	})
}

func TestParseVersions_table_not_found(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		_, err := ParseVersions(content)
		if err == nil {
			t.Errorf("ParseVersions() = nil, want error")
		}

		if err.Error() != "table not found" {
			t.Errorf("ParseVersions() = %v, want %v", err.Error(), "table not found")
		}
	})
}

func TestParseVersions_invalid_content(t *testing.T) {
	content := `<table>`

	t.Run("success", func(t *testing.T) {
		_, err := ParseVersions(content)
		if err == nil {
			t.Errorf("ParseVersions() = nil, want error")
		}
	})
}
