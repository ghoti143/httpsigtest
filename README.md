# httpsigtest

## server

`go run main.go server`

This will start up a server which uses the httpsig's [NewVerifyMiddleware](https://github.com/ghoti143/httpsig/blob/03f63b55e51b5156c5b702c3efda1e39ca5c1194/httpsig.go#L107) to verify incoming messages.

## client

`go run main.go client {method}` where {method} is one of [get, post, patch, put, delete]

This will create a request and send it using an http.Client configured to use httpsig's [NewSignTransport](https://github.com/ghoti143/httpsig/blob/03f63b55e51b5156c5b702c3efda1e39ca5c1194/httpsig.go#L35) which will sign outgoing messages.

## postman collection

This collection will send signed requests to the running server.  It has tests that ensure the responses are 200 OK and contain the message "Your request has a valid signature!"

1. Run the requests in the collection

![image](https://user-images.githubusercontent.com/624146/182958261-da624cf5-d2b5-4dd5-8c9d-5c00c1298aca.png)

2. Ensure all tests pass:

![image](https://user-images.githubusercontent.com/624146/182958108-c074680a-1ddd-461b-bda9-5555fa7567c1.png)

Todo:
1. write code in the collection's Pre-request Script to generate the sig headers
