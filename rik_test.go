package rik

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	pastebin := "https://pastebin.com/raw/seCnUfhr"
	webhook := "https://discord.com/api/webhooks/0/A"
	messages := "https://discord.com/api/v9/channels/0/messages"
	token := "token"
	filename := "hello.txt"

	//Specify the client
	DefaultClient = &http.Client{}

	s, res := Get(pastebin).MustDoReadString()
	t.Logf("GET (%d): %s", res.StatusCode, s)

	s, res = Post(webhook).
		JSON(map[string]interface{}{"content": "hello json"}).
		MustDoReadString()
	t.Logf("JSON POST (%d): %s", res.StatusCode, s)

	s, res = Post(messages).
		JSON(map[string]interface{}{"content": "hello json with header"}).
		Header(map[string][]string{"authorization": {token}}).
		MustDoReadString()
	t.Logf("JSON POST with Header (%d): %s", res.StatusCode, s)

	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	s, res = Post(webhook).
		Multipart(NewMultipart().
			File(filename, file).
			Field("content", strings.NewReader("hello multiport")).
			Boundary("END_OF_PART").
			MustBuild(),
		).
		MustDoReadString()
	t.Logf("Multipart POST with Header (%d): %s", res.StatusCode, s)

	s, res = Post(webhook).
		Form(map[string]string{"content": "hello form"}).
		MustDoReadString()
	t.Logf("Form POST (%d): %s", res.StatusCode, s)
}
