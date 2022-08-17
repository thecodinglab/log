# 🪵

`github.com/thecodinglab/log` is another structured logging library for Golang.
It has been inspired by some other awesome logging libraries such
as [logrus](https://github.com/sirupsen/logrus)
and [zap](https://github.com/uber-go/zap).

I started with this library to combine my personal logging stack into a single
library which I can simply import in all of my projects. I am still heavily
working on this, so I don't recommend using it in production. However, if you
have any feedback or suggestions, feel free to open an issue or pull request.

## Usage

```shell
go get -u github.com/thecodinglab/log
```

The logging API uses the golang `context` package to pass instances of
the `log.Logger` interface through the application. By doing this, each module
can append additional meta information to the logger without having to pass a
logging instance to every sub-service it uses.

* Application &rarr; initialize logging
    * RESTful API Service &rarr; attach HTTP client information
        * Database &rarr; attach query statement
        * ...
    * ...

```go
package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/thecodinglab/log"
	"github.com/thecodinglab/log/level"
	"github.com/thecodinglab/log/capture"
	"github.com/thecodinglab/log/sinks/file"
)

func main() {
	// initialize logger to print to stdout
	logger := log.New(
		file.New(os.Stdout, file.TextFormatter{}, level.Info),
	)
	defer logger.Sync()

	// attach the logger to the global application context
	ctx := log.Attach(context.Background(), logger)

	router := http.NewServeMux()

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: loggingMiddleware(router),

		// use the previously created application context as base context for each request
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	if err := server.ListenAndServe(); err != nil {
		capture.Error(ctx, err)
	}
}

func loggingMiddleware(base http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// attach client information to the logger
		logger := log.With(req.Context(),
			"client", req.RemoteAddr,
			"request_method", req.Method,
			"request_url", req.URL.String(),
		)

		logger.Debug("client ", req.RemoteAddr, " requested resource ", req.URL)

		// pass client information down to the next handler
		ctx := log.Attach(req.Context(), logger)
		base.ServeHTTP(res, req.WithContext(ctx))
	})
}
```
