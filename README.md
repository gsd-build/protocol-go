# protocol-go

Go types for the GSD Cloud wire protocol — the message format the daemon and relay use to talk to each other over websockets.

## Install

```bash
go get github.com/gsd-build/protocol-go
```

## Usage

```go
import (
    "encoding/json"

    protocol "github.com/gsd-build/protocol-go"
)

// Decode an incoming frame:
env, err := protocol.ParseEnvelope(data)
if err != nil {
    // handle error
}
switch msg := env.Payload.(type) {
case *protocol.Hello:
    _ = msg.MachineID
case *protocol.Task:
    _ = msg.Prompt
}

// Decode an incoming frame with transport bounds:
env, err = protocol.ParseEnvelopeWithLimits(data, protocol.DefaultEnvelopeLimits())
if err != nil {
    // close or reject the frame
}

// Encode an outgoing frame - marshal the payload struct directly:
data, _ := json.Marshal(&protocol.Hello{
    Type:      protocol.MsgTypeHello,
    MachineID: "machine-1",
})
```

`protocol.ExtractBinding`, `protocol.ValidateRequestBinding`, and
`protocol.ValidateSessionBinding` expose request/session correlation checks for
relay and daemon handlers.

See [PROTOCOL.md](./PROTOCOL.md) for the wire format specification.

## License

MIT — see [LICENSE](./LICENSE).
