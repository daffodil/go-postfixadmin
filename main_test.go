
package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO end to end testing
// https://codegangsta.gitbooks.io/building-web-apps-with-go/content/testing/end_to_end/index.html
func Test_App(t *testing.T) {

	ts := httptest.NewServer( App("config.yaml") )
	defer ts.Close()

	res, err := http.Get("/api/v1/admin/domains")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := "Before...Hello World...After"

	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
