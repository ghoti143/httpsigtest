// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ghoti143/httpsig"
	"moul.io/http2curl"
)

const secret = "support-your-local-cat-bonnet-store"
const url = "http://127.0.0.1:1234/"

func server() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "Your request has a valid signature!")
	})

	middleware := httpsig.NewVerifyMiddleware(httpsig.WithHmacSha256("key1", []byte(secret)))
	http.Handle("/", middleware(h))

	log.Fatal(http.ListenAndServe("127.0.0.1:1234", http.DefaultServeMux))
}

func client_get() {
	req, _ := http.NewRequest("GET", url, nil)	
	
	client_do(req)
}

func client_post() {	
	var jsonData = []byte(`{
		"name": "morpheus",
		"job": "leader"
	}`)	

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client_do(req)	
}

func client_patch() {
	var jsonData = []byte(`{
		"name": "neo",
		"job": "the one"
	}`)	

	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client_do(req)	
}

func client_put() {
	var jsonData = []byte(`{
		"name": "trinity",
		"job": "sidekick?"
	}`)	

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client_do(req)	
}

func client_delete() {
	req, _ := http.NewRequest("DELETE", url+"1234", nil)
	client_do(req)	
}

func client_do(req *http.Request) {

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println("\ncurl without sig: \n\n", command)

	client := http.Client{
		// Wrap the transport:
		Transport: httpsig.NewSignTransport(http.DefaultTransport,
			httpsig.WithHmacSha256("key1", []byte(secret))),
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("got err: ", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("\n\n", resp.Status, " || ", string(b))	

	command1, _ := http2curl.GetCurlCommand(resp.Request)
	fmt.Println("\n\ncurl with sig:\n\n", command1)
}

func main() {
	arg := os.Args[1]

	if arg == "server" {
		server() 
	} else {
		switch os.Args[2] {
		case "post":
			client_post()
		case "patch":
			client_patch()
		case "put":
			client_put()
		case "delete":
			client_delete()
		default:
			// TODO: unknown params could be kept? hard to say.
			client_get()
		}
		
	}

}
