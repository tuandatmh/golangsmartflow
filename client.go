package tsales_smart_flow

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type Http interface {
	Get()
	Post()
	Put()
	Delete()
	PostLogin()
}

type Client struct {
	Key       string
	Token     string
	AuthToken string
	BaseUrl   string
	Http
}

func (c *Client) PostLogin(endpoint string, params map[string]string, auths map[string]string) []byte {
	return c.execute("POST", endpoint, params, auths)
}

func (c *Client) Get(endpoint string, params map[string]string) []byte {
	return c.execute("GET", endpoint, params, nil)
}

func (c *Client) Post(endpoint string, params map[string]string) []byte {
	return c.execute("POST", endpoint, params, nil)
}

func (c *Client) Put(endpoint string, params map[string]string) []byte {
	return c.execute("PUT", endpoint, params, nil)
}

func (c *Client) Delete(endpoint string, params map[string]string) []byte {
	return c.execute("DELETE", endpoint, params, nil)
}

func (c *Client) ParamInterface(method, endpoint string, params interface{}) []byte {
	if method == "Post" {
		return c.executeInterface("POST", endpoint, params)
	} else {
		return c.executeInterface("PUT", endpoint, params)
	}
}

func (c *Client) RequestByForm(method, endpoint string, params map[string]string) []byte {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	// client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if params["file"] != "" && method == "POST" {
		fieldName := "file"
		fileName := filepath.Base(params["file"])

		// Open the file
		file, err := os.Open(params["file"])
		if err != nil {
			log.Println(err)
			return []byte(``)
		}
		contentType := func() string {
			defer func() {
				_, _ = file.Seek(0, 0)
			}()

			fileData, err := ioutil.ReadAll(file)
			if err != nil {
				return "application/octet-stream"
			}

			return http.DetectContentType(fileData)
		}()

		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
		header.Set("Content-Type", contentType)
		part, _ := writer.CreatePart(header)

		_, _ = io.Copy(part, file)

		delete(params, "file")
	}

	fmt.Println("hhhh")
	fmt.Println(params)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return []byte(``)
	}

	req, requestErr := http.NewRequest(method, c.BaseUrl+endpoint, body)
	if requestErr != nil {
		log.Fatalln(requestErr)
	}

	req.Header.Add("auth_token", c.AuthToken)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []byte(``)
	}

	return c.parseBody(resp)
}

func (c *Client) buildUrl(baseUrl, endpoint string, params map[string]string) string {
	query := make([]string, len(params))
	for k := range params {
		query = append(query, k+"="+params[k])
	}
	return baseUrl + endpoint + "?" + strings.Join(query, "&")
}

func (c *Client) parseBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte(``)
	}
	return body
}

func (c *Client) execute(method, endpoint string, params map[string]string, auths map[string]string) []byte {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	var (
		req        *http.Request
		requestErr error
	)

	if method != "GET" {
		payloadBytes, err := json.Marshal(params)
		if err != nil {
			log.Println(err)
			return []byte(``)
		}

		body := bytes.NewReader(payloadBytes)
		req, requestErr = http.NewRequest(method, c.BaseUrl+endpoint, body)
		req.Header.Add("Content-Type", "application/json")
		if auths != nil {
			req.Header.Add("Key", auths["key"])
			req.Header.Add("Token", auths["token"])
		}
	} else {
		req, requestErr = http.NewRequest(method, c.buildUrl(c.BaseUrl, endpoint, params), nil)
	}
	if requestErr != nil {
		panic(requestErr)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Auth_token", c.AuthToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return []byte(``)
	}
	return c.parseBody(resp)
}

func (c *Client) executeInterface(method, endpoint string, params interface{}) []byte {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	var (
		req        *http.Request
		requestErr error
	)

	payloadBytes, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return []byte(``)
	}

	body := bytes.NewReader(payloadBytes)

	req, requestErr = http.NewRequest(method, c.BaseUrl+endpoint, body)
	req.Header.Add("Content-Type", "application/json")

	if requestErr != nil {
		panic(requestErr)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Auth_token", c.AuthToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return []byte(``)
	}
	return c.parseBody(resp)
}
