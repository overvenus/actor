# actor - A Minimum Go Actor Framework

[![GoDoc](https://godoc.org/github.com/overvenus/actor?status.svg)](https://pkg.go.dev/github.com/overvenus/actor?tab=doc)

## Features

* A minimum actor runtime.
* Run 100k+ actors concurrently in a pool of goroutine (8*GOMAXPROCS).
* Comprehensive metrics monitoring.

## Status

This package is kind of *stable*, it's currently used by [TiFlow](https://github.com/pingcap/tiflow) project in production.

New features or bug fixes are welcome!

## Examples

### Ping pong

```go
type pingpong struct {
	peer   actor.ID
	router *actor.Router[int]
}

func (p *pingpong) Poll(ctx context.Context, msgs []message.Message[int]) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	println("recv from peer", p.peer, msgs[0].Value)
	p.router.Send(p.peer, msgs[0])
	return true
}

func (p *pingpong) OnClose() {}

func main() {
	sys, router := actor.NewSystemBuilder[int]("ping-pong").Build()
	ctx := context.Background()
	sys.Start(ctx)

	a1 := &pingpong{peer: actor.ID(2), router: router}
	mb1 := actor.NewMailbox[int](actor.ID(1), 1)
	sys.Spawn(mb1, a1)

	a2 := &pingpong{peer: actor.ID(1), router: router}
	mb2 := actor.NewMailbox[int](actor.ID(2), 2)
	sys.Spawn(mb2, a2)

	// Initiate ping pong.
	router.Send(actor.ID(1), message.ValueMessage(0))

	time.Sleep(3 * time.Second)
	sys.Stop()
}
```

