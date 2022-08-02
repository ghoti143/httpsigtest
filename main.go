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

func server() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "Your request has a valid signature!")
	})

	middleware := httpsig.NewVerifyMiddleware(httpsig.WithHmacSha256("key1", []byte(secret)))
	http.Handle("/", middleware(h))

	log.Fatal(http.ListenAndServe("127.0.0.1:1234", http.DefaultServeMux))
}

func client() {
	client := http.Client{
		// Wrap the transport:
		Transport: httpsig.NewSignTransport(http.DefaultTransport,
			httpsig.WithHmacSha256("key1", []byte(secret))),
	}

	// GET request

	//resp, err := client.Get("http://127.0.0.1:1234/")

	// END GET request

	// POST request
	
	var jsonData = []byte(`{
		"name": "morpheus",
		"job": "leader"
	}`)	

	request, err := http.NewRequest("POST", "http://127.0.0.1:1234/", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	command1, _ := http2curl.GetCurlCommand(request)
	fmt.Println(command1)

	resp, err := client.Do(request)

	// END POST request


	if err != nil {
		fmt.Println("got err: ", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
	fmt.Println(resp.Status)
	fmt.Println(resp.Request.Header)

	

	// Output:
	// 200 OK
}

func main() {
	arg := os.Args[1]

	if arg == "server" {
		server() 
	} else {
		client()
	}

}
