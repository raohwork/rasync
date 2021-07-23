some helpers to write common routines

[![GoDoc](https://godoc.org/github.com/raohwork/routines?status.svg)](https://godoc.org/github.com/raohwork/routines)
[![Go Report Card](https://goreportcard.com/badge/github.com/raohwork/routines)](https://goreportcard.com/report/github.com/raohwork/routines)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-97%25-brightgreen.svg?longCache=true&style=flat)</a>

Helpers here can:

- Prevent your crawler from getting banned (`RunAtleast(duration, task)`)
- Running task repeatly in background (`InfiniteLoop(task)`)
- Retry task until first successful attempt (`Retry(task)`)

and more.

# Race conditions

Values returned in this library are thread-safe. However, thread-safety of external
function is not covered.

Considering this example:

```go
f := Recorded(yourFunc)
```

`f` is thread-safe iff `yourFunc` is thread-safe.

### simple ws client as example

with gorilla/websocket

```go
type ws struct {
    conn *websocket.Conn
    Ctrl routines.InfiniteLoopControl
}

func (w *ws) update() error {
    typ, buf, err := w.conn.ReadMessage()
    
    if err != nil {
        return err
    }
    
    switch typ {
    case websocket.PongMessage:
        log.Print("got Pong")
    case websocket.PingMessage:
        log.Print("got Ping")
    case websocket.TextMessage:
        log.Print("got Text: ", string(buf))
    default:
        log.Print("get Message: ", buf)
    }
    
    return nil
}

func (w *ws) ping() error {
    return w.conn.WriteControl(
		websocket.PingMessage,
		[]byte(`test`),
		time.Now().Add(time.Second),
	)
}

func New(url string) (ret *ws, err error) {
    conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
    if err != nil {
        return
    }
    
    ret = &ws{ conn: conn }
    ret.Ctrl = routines.AnyErr(
        // send ping frame every 30s
        routines.InfiniteLoop(routines.RunAtLeast(30 * time.Second, ret.ping)),
        // handle messages
        routines.InfiniteLoop(ret.update),
    )
    return
}

func main() {
    cl, err := New("wss://example.com")
    if err != nil {
        log.Fatal("cannot connect: ", err)
    }
    defer cl.Ctrl.Cancel()
    
    for e := range cl.Ctrl.Err {
        log.Print("catched an error: ", e)
    }
}
```

# License

Copyright Chung-Ping Jen <ronmi.ren@gmail.com> 2021-

MPL v2.0
