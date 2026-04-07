# protocol-go

Go types for the GSD Cloud wire protocol — the message format the daemon and relay use to talk to each other over websockets.

## Install

```bash
go get github.com/gsd-build/protocol-go
```

## Usage

```go
import protocol "github.com/gsd-build/protocol-go"

env := protocol.Envelope{Type: protocol.MessageTypeHello, Payload: hello}
```

See [PROTOCOL.md](./PROTOCOL.md) for the wire format specification.

## License

MIT — see [LICENSE](./LICENSE).
