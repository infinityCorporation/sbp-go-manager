package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func operationsCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request... server operational\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Server Operational\n")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()

	statType := url.Get("type")

	switch statType {
	case "home":
		io.WriteString(w, "Home stats... Delivered.\n")
		break
	case "nfl":
		io.WriteString(w, "NFL stats... Delivered.\n")
		break
	default:
		io.WriteString(w, "No Type Provided... Defaulting to Home... Delivered.\n")
	}

}

func getInfo(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	user := url.Get("id")

	response := fmt.Sprintf("The user id is: %s\n", user)
	io.WriteString(w, response)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", operationsCheck)
	mux.HandleFunc("/stats", getStats)
	mux.HandleFunc("/user", getInfo)

	ctx, cancelCtx := context.WithCancel(context.Background())

	serverOne := &http.Server{
		Addr:    ":1000",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":2000",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server two: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
