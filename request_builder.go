package rik

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	ContentTypeJSON        = "application/json"
	ContentTypeForm        = "application/x-www-form-urlencoded"
	ContentTypeOctetStream = "application/octet-stream"
	ContentTypeText        = "text/plain"
)

type RequestBuilder struct {
	urlData     *url.URL
	method      string
	headerData  *http.Header
	contentType string
	bodyData    io.Reader
	queryData   *url.Values
	formData    *url.Values
	client      *http.Client
}

func NewRequest(urlData, method string) *RequestBuilder {
	u, err := url.Parse(urlData)
	if err != nil {
		panic(err)
	}
	return &RequestBuilder{
		urlData:    u,
		method:     method,
		headerData: &http.Header{},
		queryData:  &url.Values{},
		formData:   &url.Values{},
	}
}

func (r *RequestBuilder) Client(client *http.Client) *RequestBuilder {
	r.client = client
	return r
}

func (r *RequestBuilder) Header(key string, value ...string) *RequestBuilder {
	for _, v := range value {
		r.headerData.Add(key, v)
	}
	return r
}

func (r *RequestBuilder) Headers(data http.Header) *RequestBuilder {
	for k := range data {
		r.Header(k, data.Values(k)...)
	}
	return r
}

func (r *RequestBuilder) Query(key string, value ...string) *RequestBuilder {
	for _, v := range value {
		r.queryData.Add(key, v)
	}
	r.urlData.RawQuery = r.queryData.Encode()
	return r
}

func (r *RequestBuilder) Queries(data url.Values) *RequestBuilder {
	for k := range data {
		r.queryData.Add(k, data.Get(k))
	}
	r.urlData.RawQuery = r.queryData.Encode()
	return r
}

func (r *RequestBuilder) Body(contentType string, data io.Reader) *RequestBuilder {
	r.contentType = contentType
	r.bodyData = data
	return r
}

func (r *RequestBuilder) JSON(data any) *RequestBuilder {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return r.Body(ContentTypeJSON, bytes.NewBuffer(b))
}

func (r *RequestBuilder) Form(key string, value ...string) *RequestBuilder {
	for _, v := range value {
		r.formData.Add(key, v)
	}
	return r.Body(ContentTypeForm, strings.NewReader(r.formData.Encode()))
}

func (r *RequestBuilder) Forms(data url.Values) *RequestBuilder {
	for k, v := range data {
		r.Form(k, v...)
	}
	return r.Body(ContentTypeForm, strings.NewReader(r.formData.Encode()))
}

func (r *RequestBuilder) Text(data string) *RequestBuilder {
	return r.Body(ContentTypeText, strings.NewReader(data))
}

func (r *RequestBuilder) File(data *os.File) *RequestBuilder {
	return r.Body(ContentTypeOctetStream, data)
}

func (r *RequestBuilder) Multipart(data *MultipartData) *RequestBuilder {
	return r.Body(data.ContentType, data.Buffer)
}

func (r *RequestBuilder) Build() (*http.Request, error) {
	req, err := http.NewRequest(r.method, r.urlData.String(), r.bodyData)
	if err == nil {
		if r.headerData != nil {
			req.Header = *r.headerData
		}
		if req.Header.Get("Content-Type") == "" && r.contentType != "" {
			req.Header.Set("Content-Type", r.contentType)
		}
	}
	return req, err
}

func (r *RequestBuilder) MustBuild() *http.Request {
	req, err := r.Build()
	if err != nil {
		panic(err)
	}
	return req
}

func (r *RequestBuilder) DoClient(client *http.Client) (*http.Response, error) {
	var c *http.Client
	if DefaultClient != nil {
		c = DefaultClient
	}
	if client != nil {
		c = client
	}
	if c == nil {
		return nil, errors.New("no client specified")
	}
	req, err := r.Build()
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *RequestBuilder) Do() (*http.Response, error) {
	return r.DoClient(r.client)
}

func (r *RequestBuilder) MustDo() *http.Response {
	res, err := r.Do()
	if err != nil {
		panic(err)
	}
	return res
}

func (r *RequestBuilder) DoReadByteClient(client *http.Client) ([]byte, *http.Response, error) {
	res, err := r.DoClient(client)
	if err != nil {
		return nil, res, err
	}
	b, err := ReadByte(res)
	if err != nil {
		return nil, res, err
	}
	return b, res, nil
}

func (r *RequestBuilder) DoReadByte() ([]byte, *http.Response, error) {
	return r.DoReadByteClient(r.client)
}

func (r *RequestBuilder) DoReadStringClient(client *http.Client) (string, *http.Response, error) {
	b, res, err := r.DoReadByteClient(client)
	if err != nil {
		return "", res, err
	}
	return string(b), res, nil
}

func (r *RequestBuilder) DoReadString() (string, *http.Response, error) {
	return r.DoReadStringClient(r.client)
}

func (r *RequestBuilder) MustDoReadByteClient(client *http.Client) ([]byte, *http.Response) {
	b, res, err := r.DoReadByteClient(client)
	if err != nil {
		panic(err)
	}
	return b, res
}

func (r *RequestBuilder) MustDoReadByte() ([]byte, *http.Response) {
	return r.MustDoReadByteClient(r.client)
}

func (r *RequestBuilder) MustDoReadStringClient(client *http.Client) (string, *http.Response) {
	b, res := r.MustDoReadByteClient(client)
	return string(b), res
}

func (r *RequestBuilder) MustDoReadString() (string, *http.Response) {
	return r.MustDoReadStringClient(r.client)
}
