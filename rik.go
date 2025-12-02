package rik

import (
	"io"
	"net/http"
)

var DefaultClient *http.Client

func Get(url string) *RequestBuilder {
	return NewRequest(url, "GET")
}

func Post(url string) *RequestBuilder {
	return NewRequest(url, "POST")
}

func Put(url string) *RequestBuilder {
	return NewRequest(url, "PUT")
}

func Patch(url string) *RequestBuilder {
	return NewRequest(url, "PATCH")
}

func Delete(url string) *RequestBuilder {
	return NewRequest(url, "DELETE")
}

func Options(url string) *RequestBuilder {
	return NewRequest(url, "OPTIONS")
}

func Head(url string) *RequestBuilder {
	return NewRequest(url, "HEAD")
}

func Trace(url string) *RequestBuilder {
	return NewRequest(url, "TRACE")
}

func Connect(url string) *RequestBuilder {
	return NewRequest(url, "CONNECT")
}

func ReadByte(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReadString(res *http.Response) (string, error) {
	b, err := ReadByte(res)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func MustReadByte(res *http.Response) []byte {
	b, err := ReadByte(res)
	if err != nil {
		panic(err)
	}
	return b
}

func MustReadString(res *http.Response) string {
	return string(MustReadByte(res))
}
