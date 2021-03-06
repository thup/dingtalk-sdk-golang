package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var httpClient = &http.Client{}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
}

func Post(url string, params map[string]string, body string) (string, error) {
	resp, err := httpClient.Post(url+ConvertToQueryParams(params), "application/json", strings.NewReader(body))
	return ResponseHandle(resp, err)
}

func PostFile(url string, params map[string]string, path string, name string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return "", err
	}
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	_, _ = io.Copy(fileWriter, fh)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := httpClient.Post(url+ConvertToQueryParams(params), contentType, bodyBuf)
	return ResponseHandle(resp, err)
}

func PostFileWithReader(url string, params map[string]string, reader io.Reader) (string, error) {
	resp, err := httpClient.Post(url+ConvertToQueryParams(params), "multipart/form-data", reader)
	return ResponseHandle(resp, err)
}

func Get(url string, params map[string]string) (string, error) {
	resp, err := httpClient.Get(url + ConvertToQueryParams(params))
	return ResponseHandle(resp, err)
}

func ResponseHandle(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ConvertToQueryParams(params map[string]string) string {
	if &params == nil || len(params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		buffer.WriteString(k + "=" + v + "&")
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}
