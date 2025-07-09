package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const keyServerAddr MyString = "serverAddr"

type MyString string

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	server := &http.Server{
		Addr:    "127.0.0.1:3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server one: %s\n", err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s got / request first(%t)=%s, second(%t)=%s body:\n%s\n",
		context.Value(keyServerAddr),
		hasFirst,
		first,
		hasSecond,
		second,
		bodyData,
	)
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	fmt.Printf("%s got /hello request\n", context.Value(keyServerAddr))
	io.WriteString(w, "Hello, Golang!\n")
}
