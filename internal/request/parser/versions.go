package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func versionsTable(n *html.Node, class string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == class {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := versionsTable(c, class); found != nil {
			return found
		}
	}

	return nil
}

func versionRow(tr *html.Node) (string, string) {
	var version string
	var label string

	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			version = extractVersion(c)
			break
		}
	}

	if version == "" {
		return "", ""
	}

	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			label = extractLabel(c)
			if label != "" {
				break
			}
		}
	}

	return version, label
}

func extractVersion(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if text := extractVersion(c); text != "" {
			return text
		}
	}

	return ""
}

func extractLabel(c *html.Node) string {
	var label string

	for inp := c.FirstChild; inp != nil; inp = inp.NextSibling {
		if inp.Type != html.ElementNode || inp.Data != "input" {
			continue
		}

		var valueVal string
		isChecked := false

		for _, attr := range inp.Attr {
			if attr.Key == "checked" {
				isChecked = true
				break
			}
		}

		if isChecked {
			for _, attr := range inp.Attr {
				if attr.Key == "value" {
					valueVal = attr.Val
				}
			}
		}

		if valueVal != "" {
			label = valueVal
			break
		}
	}

	return label
}
