# goshmuffle
The simple way to start shell commands from golang.
Also, it can store stdout to the provided interface.
Use it on your own risk :)

## Usage

### Get the go-lib module

```bash
go get github.com/pershinov/goshmuffle@v1.0.0
```

### Example
```go
package main

import (
	"context"
	"fmt"

	"github.com/pershinov/goshmuffle"
)

type res struct {
	s string
}

func (r *res) Store(s string) {
	r.s = s
}

func main() {
	r := &res{}
	cmd := goshmuffle.New("echo", "hello").WithResult(r)

	err := cmd.Run(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r.s)
}
```

### Example with terminate
```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pershinov/goshmuffle"
)

type result struct {
	s []string
}

func (r *result) Store(s string) {
	fmt.Println(s)
	r.s = append(r.s, s)
}

func main() {
	r := &result{}
	cmd := goshmuffle.New("ping", "8.8.8.8").WithResult(r)

	go func() {
		err := cmd.Run(context.Background())
		if err != nil && !strings.Contains(err.Error(), "terminated") {
			fmt.Println(err)
		}
	}()

	// waiting for running
	for !cmd.IsRunning() {
	}

	time.Sleep(5 * time.Second)
	err := cmd.Terminate()
	if err != nil {
		fmt.Println(err)
	}

	// waiting for done
	for !cmd.IsDone() {
	}

	fmt.Println(r.s)
	fmt.Println(cmd.IsDone(), cmd.IsRunning())
}
```

## Have a good vibe ^..^