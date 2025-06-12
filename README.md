# dify-sdk-go

This is the Go SDK for the Dify API, which allows you to easily integrate Dify into your Go applications.

## Install

```
go get github.com/taadis/dify-sdk-go
```

## Usage

After installing the SDK, you can use it in your go project like this:

```go
package main

import (
	"context"
	"log"
	"strings"

	"github.com/taadis/dify-sdk-go"
)

func main() {
	ctx = context.Background()
	client = dify.NewClient("your-dify-api-host", "your-api-key")
	req := &dify.ChatMessageRequest{
		Query: "your-question",
		User: "your-user",
	}

	var	ch chan dify.ChatMessageStreamChannelResponse
	var err error
	ch, err = c.Api().ChatMessagesStream(ctx, req); err != nil {
		return
	}

	var strBuilder strings.Builder

	for {
		select {
		case <-ctx.Done():
			return
		case streamData, isOpen := <-ch:
			if err = streamData.Err; err != nil {
				log.Println(err.Error())
				return
			}
			if !isOpen {
				log.Println(strBuilder.String())
				return
			}

			strBuilder.WriteString(streamData.Answer)
		}
	}
}

```

## License

This SDK is released under the MIT License.
