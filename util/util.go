package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Call(result interface{}, opts ReqOpts, debugLog bool) error {
	targetURL := fmt.Sprintf("%s%s", opts.Host, opts.RelativeURL)

	req, err := http.NewRequest(opts.Method, targetURL, bytes.NewReader(opts.Body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	if debugLog {
		log.Printf("REQUEST: [%s] %s - %+v - %+v\n", req.Method, req.URL.String(), req.Header, string(opts.Body))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("result is not defined")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	if debugLog {
		log.Printf("RESPONSE: %+v\n", result)
	}

	return err
}

func RunServerGracefully(port, timeoutGracefull int, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutGracefull)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
