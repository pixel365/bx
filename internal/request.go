package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	jar        http.CookieJar
}

func NewClient(client HTTPClient, jar http.CookieJar) *Client {
	return &Client{
		httpClient: client,
		jar:        jar,
	}
}

// Authorization performs user authentication by sending login credentials
// to the Bitrix Partner Portal.
//
// It sends a POST request with the login and password as form data and checks
// the response for authentication success by verifying the presence of a
// "BITRIX_SM_LOGIN" cookie.
//
// Returns a slice of cookies if authentication is successful or an error if
// authentication fails or an issue occurs during the request.
func (c *Client) Authorization(login, password string) ([]*http.Cookie, error) {
	if login == "" {
		return nil, errors.New("empty login")
	}

	if password == "" {
		return nil, errors.New("empty password")
	}

	body := url.Values{
		"AUTH_FORM":     {"Y"},
		"TYPE":          {"AUTH"},
		"USER_LOGIN":    {login},
		"USER_PASSWORD": {password},
		"USER_REMEMBER": {"Y"},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://partners.1c-bitrix.ru/personal/",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var cookies []*http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "BITRIX_SM_LOGIN" && c.Value == login {
			cookies = resp.Cookies()
			break
		}
	}

	if len(cookies) == 0 {
		return nil, errors.New("authentication failed")
	}

	return cookies, nil
}

// UploadZIP uploads a ZIP file containing the module's data to the Bitrix Partner Portal.
//
// This function first validates that the module and cookies are provided. It then retrieves the
// session ID and prepares the ZIP file for upload. The request is sent with the necessary form
// data, including the session ID, module name, and the ZIP file. The response body is checked
// for the result of the upload operation.
//
// Parameters:
// - module: The module whose ZIP file is being uploaded.
// - cookies: The cookies containing the authentication information.
//
// Returns:
// - An error if any step fails (e.g., missing session, file errors, upload failure).
func (c *Client) UploadZIP(module *Module, cookies []*http.Cookie) error {
	if module == nil {
		return errors.New("module is nil")
	}

	if cookies == nil {
		return errors.New("cookies is nil")
	}

	session := c.SessionId(module, cookies)
	if session == "" {
		return errors.New("no session")
	}

	path, err := module.ZipPath()
	if err != nil {
		return err
	}

	path = filepath.Clean(path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fileWriter, err := writer.CreateFormFile("update", module.Version+".zip")
	if err != nil {
		return err
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		return err
	}

	_ = writer.WriteField("sessid", session)
	_ = writer.WriteField("ID", module.Name)
	_ = writer.WriteField("submit", "Y")

	err = writer.Close()
	if err != nil {
		return err
	}

	u, _ := url.Parse("https://partners.1c-bitrix.ru/personal/modules/deploy.php")
	req, err := http.NewRequestWithContext(module.Ctx, http.MethodPost, u.String(), &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	c.jar.SetCookies(u, cookies)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}(resp)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := uploadResult(string(respBody)); err != nil {
		return err
	}

	return nil
}

// SessionId retrieves the session ID for a given module from the Bitrix Partner Portal.
//
// The function sends a GET request to the edit page of the module, then parses the HTML
// response to extract the session ID. The session ID is needed for later operations
// like uploading data to the portal.
//
// Parameters:
// - module: The module for which the session ID is being retrieved.
// - cookies: The cookies containing the authentication information.
//
// Returns:
// - The session ID as a string if found, otherwise returns an empty string.
func (c *Client) SessionId(module *Module, cookies []*http.Cookie) string {
	u, _ := url.Parse("https://partners.1c-bitrix.ru/personal/modules/edit.php?ID=" + module.Name)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return ""
	}

	c.jar.SetCookies(u, cookies)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}(resp)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	doc, err := html.Parse(strings.NewReader(string(respBody)))
	if err != nil {
		return ""
	}

	session := ""
	sid := "sessid"
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var name, value string
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == sid {
					name = attr.Val
				}
				if attr.Key == "value" {
					value = attr.Val
				}
			}
			if name == sid {
				session = value
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return session
}

// uploadResult processes the HTML content returned from the upload request
// to check for error messages.
//
// The function parses the HTML content and searches for a <p> element with
// a specific CSS class (`paragraph-15 color-red m-0`), which indicates
// an error message. If such an element is found, the error message is
// extracted and returned as an error.
//
// Parameters:
// - htmlContent: The HTML response body to be parsed for error messages.
//
// Returns:
//   - An error if an error message is found in the HTML content or nil if
//     no errors are present.
func uploadResult(htmlContent string) error {
	var err error

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			hasClass := false
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "paragraph-15 color-red m-0" {
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
