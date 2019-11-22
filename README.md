# :handshake: handgover

</p>
<p align="left">
	<a alt="Licence" >
		<img src="https://img.shields.io/github/license/NewStore-oss/handgover">
	</a>
	<a alt="Build" >
		<img src="https://github.com/newstore-oss/handgover/workflows/Go/badge.svg">
	</a>
	<a alt="Go Module" >
		<img src="https://img.shields.io/github/go-mod/go-version/newstore-oss/handgover">
	</a>
	<a alt="Stars" >
		<img src="https://img.shields.io/github/stars/Newstore-oss/handgover">
	</a>
</p>

handgover is a tool to fill your struct based on your own defined tags and matching sources.


## Overview
handgover is written in go. To analyse the given structs it uses the standard [`reflect`](https://golang.org/pkg/reflect/) package. No other third-party package is needed.


### Supported Types
 - string
 - integer (int8, int16, int32, int64, Uint, Uint8, Uint16, UInt32, UInt64)
 - Bool
 - float (float32, float64)
 - time.Duration
 - time.Time (RFC3339)
 - []byte

> **Note**: Every listed type supports *pointer* and *slice* as well.

## Usage

### Define sources
```go
sources := []handgover.Source{
    {
	Tag: "query",
	Get: func(field string) ([]string, error) {
	    return req.URL.Query()[field], nil
	},
    },
}
```

### Define your struct
```go
type MyStruct struct {
    Count int `query:"count"`
    Offset int `query:"offset"`
    Query string `query:"q"`
}
```
> **Note**:  Multiple tags per property are supported.  Source values are taken out of the order as you defined in your sources.

### Putting everything together

```go
package main

import (
	"log"
	"net/http"
	"github.com/newstore/handgover"
)

type  MyRequest  struct {
	Count int  `query:"count"`
	Offset int  `query:"offset"`
	Query string  `query:"q"`
}

func main() {
	incomingReq, _ := http.NewRequest(
		"GET",
		"http://www.example.com/?count=100&offset=abc&q=test",
		nil,
	)

	var myRequest MyRequest
	if err := Pick(incomingReq, &myRequest); err !=  nil {
		if hErr, ok := handgover.FromError(err); ok; {
			log.Println(hErr.Field)
			// OUTPUT: offset
			log.Println(hErr.Source)
			// OUTPUT: query
			log.Println(hErr.Value)
			// OUTPUT: abc
			log.Println(hErr.InnerError)
			// OUTPUT: strconv.ParseInt: parsing "abc": invalid syntax
		}
		log.Fatal(err)
		// OUTPUT: failed to set field "offset" from source "query":
		// strconv.ParseInt: parsing "abc": invalid syntax
	}

	log.Printf("%+v", myRequest)
	// OUTPUT: {Count:100 Offset:200 Query:test}
}

func Pick(req *http.Request, v interface{}) error {
	sources := []handgover.Source{
		{
			Tag: "query",
			Get: func(field string) ([]string, error) {
				return req.URL.Query()[field], nil
			},
		},
	}
	return handgover.From(sources).To(v)
}
```

## Motivation
When you create a new HTTP endpoint you probably need to get some values from your query. This is done in Go pretty well and you can easy achieve it. But when you look at the return value - it's a *string*.  It turned out that the real world is a little bit different and you may need it as specific type e.g. *integer*.  In one place you want to check it against some condition (e.g. `if count>100`) or forward it to your next component which only accepts a specific type.

Next steps would be to parse each single query parameter to your specific value and do of course error handling, because  someone is always using your API wrong `(e.g. count=abc)`. Means you have to take care about it as well. Doing that again and again for several endpoint felt tedious.

At this point the question came up "Is there no easier way of doing that?" - The idea of handgover was born :hatching_chick:!

## Licence
MIT [License](https://github.com/NewStore-oss/handgover/blob/master/LICENSE) Copyright (c) 2019 NewStore GmbH
