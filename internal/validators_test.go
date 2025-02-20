package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateModuleName_NotExisting(t *testing.T) {
	t.Run("TestValidateModuleName_NotExisting", func(t *testing.T) {
		if err := ValidateModuleName("not_exists", "./"); err != nil {
			t.Error(err)
		}
	})
}

func TestValidateModuleName_Existing(t *testing.T) {
	t.Run("TestValidateModuleName_Existing", func(t *testing.T) {
		name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
		filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", ".", name))
		if err != nil {
			t.Error()
		}

		err = os.WriteFile(filePath, []byte(DefaultYAML()), 0600)
		if err != nil {
			t.Error(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Error(err)
			}
		}(filePath)

		err = ValidateModuleName(name, ".")
		if err == nil {
			t.Errorf("error expected")
		}
	})
}
