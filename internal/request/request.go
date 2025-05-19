package request

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pixel365/bx/internal/client"
	"github.com/pixel365/bx/internal/request/parser"
	"github.com/pixel365/bx/internal/types"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
)

var getSessionFunc = getSession

// Authenticate performs user authentication by sending login credentials
// to the Bitrix Partner Portal.
//
// It sends a POST request with the login and password as form data and checks
// the response for authentication success by verifying the presence of a
// "BITRIX_SM_LOGIN" cookie.
//
// Returns a slice of cookies if authentication is successful or an error if
// authentication fails or an issue occurs during the request.
func Authenticate(client client.HTTPClient, login, password string) ([]*http.Cookie, error) {
	if login == "" {
		return nil, errors2.ErrEmptyLogin
	}

	if password == "" {
		return nil, errors2.ErrEmptyPassword
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

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer helpers.Cleanup(resp.Body, nil)

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
		return nil, errors2.ErrAuthentication
	}

	return cookies, nil
}

func Versions(
	ctx context.Context,
	client client.HTTPClient,
	module *module.Module,
	cookies []*http.Cookie,
) (types.Versions, error) {
	if module == nil {
		return nil, errors2.ErrNilModule
	}

	if cookies == nil {
		return nil, errors2.ErrNilCookie
	}

	session := sessionId(client, module, cookies)
	if session == "" {
		return nil, errors2.ErrEmptySession
	}

	u, _ := url.Parse("https://partners.1c-bitrix.ru/personal/modules/update.php?ID=" + module.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	client.SetCookies(u, cookies)

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer helpers.Cleanup(resp.Body, nil)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parser.ParseVersions(string(respBody))
}

func ChangeLabels(
	client client.HTTPClient,
	module *module.Module,
	cookies []*http.Cookie,
	versions types.Versions,
) error {
	if module == nil {
		return errors2.ErrNilModule
	}

	if cookies == nil {
		return errors2.ErrNilCookie
	}

	session := sessionId(client, module, cookies)
	if session == "" {
		return errors2.ErrEmptySession
	}

	body := url.Values{
		"sessid": {session},
		"ID":     {module.Name},
		"submit": {"Y"},
	}
	for version, label := range versions {
		body.Set(version, string(label))
	}

	u, _ := url.Parse("https://partners.1c-bitrix.ru/personal/modules/update.php")
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client.SetCookies(u, cookies)

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer helpers.Cleanup(resp.Body, nil)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return parser.UploadResult(string(respBody))
}

// UploadZIP uploads a ZIP file containing the module's data to the Bitrix Partner Portal.
//
// This function first validates that the module and cookies are provided. It then retrieves the
// session ID and prepares the ZIP file for upload. The request is sent with the necessary form
// data, including the session ID, module name, and the ZIP file. The response body is checked
// for the result of the upload operation.
//
// Parameters:
//   - ctx: context.Context.
//   - module: The module whose ZIP file is being uploaded.
//   - cookies: The cookies containing the authentication information.
//
// Returns:
//   - An error if any step fails (e.g., missing session, file errors, upload failure).
func UploadZIP(
	ctx context.Context,
	client client.HTTPClient,
	module *module.Module,
	cookies []*http.Cookie,
) error {
	if module == nil {
		return errors2.ErrNilModule
	}

	if cookies == nil {
		return errors2.ErrNilCookie
	}

	session := sessionId(client, module, cookies)
	if session == "" {
		return errors2.ErrEmptySession
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
	defer helpers.Cleanup(file, nil)

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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client.SetCookies(u, cookies)

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer helpers.Cleanup(resp.Body, nil)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return parser.UploadResult(string(respBody))
}

// sessionId retrieves the session ID for a given module from the Bitrix Partner Portal.
//
// The function sends a GET request to the edit page of the module, then parses the HTML
// response to extract the session ID. The session ID is needed for later operations
// like uploading data to the portal.
//
// Parameters:
//   - module: The module for which the session ID is being retrieved.
//   - cookies: The cookies containing the authentication information.
//
// Returns:
//   - The session ID as a string if found, otherwise returns an empty string.
func sessionId(client client.HTTPClient, module *module.Module, cookies []*http.Cookie) string {
	return getSessionFunc(client, module, cookies)
}

func getSession(client client.HTTPClient, module *module.Module, cookies []*http.Cookie) string {
	if module == nil || len(cookies) == 0 || module.Name == "" {
		return ""
	}

	u, _ := url.Parse("https://partners.1c-bitrix.ru/personal/modules/edit.php?ID=" + module.Name)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return ""
	}

	client.SetCookies(u, cookies)

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer helpers.Cleanup(resp.Body, nil)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return parser.ParseSessionId(string(respBody))
}
