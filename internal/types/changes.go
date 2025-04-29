package types

import "strings"

type Changes struct {
	Added    []string
	Modified []string
	Deleted  []string
	Moved    []string
}

// IsChangedFile checks whether the given file path corresponds to a file that has been added or modified.
//
// Parameters:
//   - path: The file path to check.
//
// Returns:
//   - true if the file is in the list of added or modified files.
//   - false otherwise.
//
// Behavior:
//   - Iterates through the `Added` and `Modified` slices of the `Changes` struct.
//   - Uses `strings.HasSuffix` to check if the provided path ends with any of the added or modified file names.
//   - Returns true on the first match, otherwise false.
//
// Notes:
//   - This method assumes that `o.Added` and `o.Modified` contain relative or base file names.
//   - If files are stored with full paths in `o.Added` and `o.Modified`, this method may produce false negatives.
//
// Example:
//
//	changes := Changes{
//	    Added:    []string{"file1.txt", "dir/file2.go"},
//	    Modified: []string{"config.yaml"},
//	}
//	fmt.Println(changes.IsChangedFile("project/dir/file2.go")) // true
//	fmt.Println(changes.IsChangedFile("config.yaml"))         // true
//	fmt.Println(changes.IsChangedFile("untracked.txt"))       // false
func (o *Changes) IsChangedFile(path string) bool {
	for _, f := range o.Added {
		if strings.HasSuffix(path, f) {
			return true
		}
	}

	for _, f := range o.Modified {
		if strings.HasSuffix(path, f) {
			return true
		}
	}

	for _, f := range o.Moved {
		if strings.HasSuffix(path, f) {
			return true
		}
	}

	return false
}
