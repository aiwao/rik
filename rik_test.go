package rik

import (
    "net/http"
    "net/url"
    "os"
    "strings"
    "testing"
)

func logResponse(t *testing.T, res *http.Response, s string) {
    t.Logf("Status %d", res.StatusCode)
    t.Log(s)
}

func Test(t *testing.T) {
    //https://github.com/aiwao/request-tester
    requestTester := "https://request-tester.cloudflare2025-5fc.workers.dev/"
    client := &http.Client{}
    //Test Get request
    {
        s, res, err := Get(requestTester).
            //Set http client
            Client(client).
            DoReadString()
        if err != nil {
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
    //Test Post request with JSON and JSONBuilder
    {
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
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
    //Set default http client for all requests
    DefaultClient = client
    //Test Put request with FormData
    {
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
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
    //Test Post request with Query parameters
    {
        s, res, err := Post(requestTester).
            Query("Key", "Value1").
            Query("Key", "Value2").
            Query("Key2", "Key2Value1", "Key2Value2").
            Queries(url.Values{
                "Key3": {"Key3Value1", "Key3Value2"},
                "Key4": {"Key4Value1", "Key4Value2"},
            }).
            DoReadString()
        if err != nil {
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
    //Test Post request with Headers and File
    {
        file, err := os.Open("hello.txt")
        if err != nil {
            t.Log(err)
            return
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
            DoReadString()
        if err != nil {
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
    //Test Post request with MultipartData and MultipartBuilder
    {
        file, err := os.Open("hello.txt")
        if err != nil {
            t.Log(err)
            return
        }
        defer file.Close()
        s, res, err := Post(requestTester).
            Multipart(NewMultipart().
                File("hello.txt", file).
                Field("content", strings.NewReader("data")).
                Boundary("END_OF_PART").
                MustBuild(),
            ).
            DoReadString()
        if err != nil {
            t.Log(err)
            return
        }
        logResponse(t, res, s)
    }
}
