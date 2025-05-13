package parser

import (
	"errors"
	"strings"

	"golang.org/x/net/html"

	"github.com/pixel365/bx/internal/types"
)

const (
	input = "input"
	value = "value"
	class = "class"
)

// ParseVersions extracts version flags from the provided HTML content.
//
// Parameters:
//   - content (string): A string containing HTML content that includes a table of version rows.
//
// Returns:
//   - Versions: A map where each key is a version string (e.g., "3.0.10")
//     and the value is a VersionLabel constant (Alpha, Beta, or Stable) representing the selected flag.
//   - error: An error if the target versions table is not found in the HTML; otherwise, nil.
//
// Description:
// ParseVersions parses the given HTML content and looks for a <table> element with the class "data-table mt-3 mb-3".
// It iterates over each <tr> row within the table's <tbody> section,
// extracting the version string (from the first <td>)
// and identifying which radio input is currently checked (determining the selected VersionLabel).
// The results are collected into a Versions map, where keys are version identifiers and values are the selected flags.
// If the expected table is not found, the function returns an error.
func ParseVersions(content string) (types.Versions, error) {
	var err error

	doc, _ := html.Parse(strings.NewReader(content))

	versions := make(map[string]types.VersionLabel)

	table := versionsTable(doc, "data-table mt-3 mb-3")
	if table == nil {
		return nil, errors.New("table not found")
	}

	for tr := table.FirstChild; tr != nil; tr = tr.NextSibling {
		if tr.Type == html.ElementNode && tr.Data == "tbody" {
			for row := tr.FirstChild; row != nil; row = row.NextSibling {
				if row.Type == html.ElementNode && row.Data == "tr" {
					version, selected := versionRow(row)
					if version != "" {
						versions[version] = types.VersionLabel(selected)
					}
				}
			}
		}
	}

	return versions, err
}

// ParseSessionId extracts the sessid value from the provided HTML content.
//
// Parameters:
//   - content (string): A string containing the HTML content to search for an input element with name="sessid".
//
// Returns:
//   - string: The value of the "value" attribute from the found input element with name="sessid".
//     If no such element is found, it returns an empty string.
//
// Description:
// The function parses the given HTML content and recursively traverses the DOM tree
// to locate an <input> element where the name attribute equals "sessid".
// Once found, it extracts and returns the value of its value attribute.
func ParseSessionId(content string) string {
	doc, _ := html.Parse(strings.NewReader(content))

	session := ""
	sid := "sessid"
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == input {
			var name, val string
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == sid {
					name = attr.Val
				}
				if attr.Key == value {
					val = attr.Val
				}
			}
			if name == sid {
				session = val
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return session
}

// UploadResult processes the HTML content returned from the upload request
// to check for error messages.
//
// The function parses the HTML content and searches for a <p> element with
// a specific CSS class (`paragraph-15 color-red m-0`), which indicates
// an error message. If such an element is found, the error message is
// extracted and returned as an error.
//
// Parameters:
//   - htmlContent: The HTML response body to be parsed for error messages.
//
// Returns:
//   - An error if an error message is found in the HTML content or nil if
//     no errors are present.
func UploadResult(content string) error {
	var err error

	doc, _ := html.Parse(strings.NewReader(content))

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			hasClass := false
			for _, attr := range n.Attr {
				if attr.Key == class && attr.Val == "paragraph-15 color-red m-0" {
					hasClass = true
					break
				}
			}
			if hasClass && n.FirstChild != nil {
				err = errors.New(strings.TrimSpace(n.FirstChild.Data))
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return err
}
