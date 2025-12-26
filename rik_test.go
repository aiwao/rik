package rik

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// https://github.com/aiwao/request-tester
const requestTester = "https://request-tester.cloudflare2025-5fc.workers.dev/"

func logResponse(t *testing.T, res *http.Response, s string) {
	t.Logf("Status %d", res.StatusCode)
	t.Log(s)
}

func TestGet(t *testing.T) {
	client := &http.Client{}
	s, res, err := Get(requestTester).
		Client(client).
		DoReadString()
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestPostJSON(t *testing.T) {
	client := &http.Client{}
	s, res, err := Post(requestTester).
		JSON(NewJSON().
			Set("Key", "Value").
			Set("Keys", []any{"Value1"}).
			Add("Keys", "Value2").
			Build(),
		).
		//Set http client and do request
		DoReadStringClient(client)
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestPut(t *testing.T) {
	client := &http.Client{}
	//Set default http client for all requests
	DefaultClient = client
	s, res, err := Put(requestTester).
		Form("Key", "Value1").
		Form("Key", "Value2").
		Form("Key2", "Key2Value1", "Key2Value2").
		Forms(url.Values{
			"Key3": {"Key3Value1", "Key3Value2"},
			"Key4": {"Key4Value1", "Key4Value2"},
		}).
		DoReadString()
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestPostQuery(t *testing.T) {
	client := &http.Client{}
	s, res, err := Post(requestTester).
		Query("Key", "Value1").
		Query("Key", "Value2").
		Query("Key2", "Key2Value1", "Key2Value2").
		Queries(url.Values{
			"Key3": {"Key3Value1", "Key3Value2"},
			"Key4": {"Key4Value1", "Key4Value2"},
		}).
		DoReadStringClient(client)
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestPostFileWithHeader(t *testing.T) {
	client := &http.Client{}
	file, err := os.Open("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	s, res, err := Post(requestTester).
		Header("Key", "Value1").
		Header("Key", "Value2").
		Header("Key2", "Value1", "Value2").
		Headers(http.Header{
			"Key3": {"Key3Value1", "Key3Value2"},
			"Key4": {"Key4Value1", "Key4Value2"},
		}).
		File(file).
		DoReadStringClient(client)
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestPostMultipart(t *testing.T) {
	client := &http.Client{}
	file, err := os.Open("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	s, res, err := Post(requestTester).
		Multipart(NewMultipart().
			File("hello.txt", file).
			Field("content", strings.NewReader("data")).
			Boundary("END_OF_PART").
			MustBuild(),
		).
		DoReadStringClient(client)
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}

func TestWithContext(t *testing.T) {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, res, err := Get(requestTester + "/timeout").
		Context(ctx).
		DoReadStringClient(client)
	if err != nil {
		t.Fatal(err)
	}
	logResponse(t, res, s)
}
