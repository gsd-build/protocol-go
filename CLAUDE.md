# CLAUDE.md

This file provides guidance to Claude Code and other AI agents when working in this repository.

## Repository Overview

`protocol-go` is the public Go module for the GSD Cloud wire protocol. It defines the JSON message envelope and typed payloads shared by the daemon and the cloud relay.

`PROTOCOL.md` is the authoritative wire-format contract. The Go structs, parsing helpers, tests, daemon consumer, and relay consumer must match that contract.

## Commands

```bash
go test ./...
go test -run TestName ./...
go test -race ./...
```

## Change Rules

- Keep this repo free of cloud-proprietary code.
- Treat message shape changes as cross-repo work. Update `PROTOCOL.md`, Go types, and tests here first.
- After protocol changes merge, tag a new module version and bump consumers in `gsd-build-daemon` and `gsd-build-cloud-app/apps/relay`.
- Preserve backward compatibility when possible. Additive fields with safe zero values are preferred.
- JSON field names are part of the wire contract. Renames and removals require coordinated daemon and relay releases.
- Tests should verify parsing, serialization, unknown/invalid payload handling, and protocol compatibility behavior.

## Release Flow

Protocol changes ship in this order:

1. Merge this repo.
2. Tag the new module version.
3. Bump `github.com/gsd-build/protocol-go` in daemon and cloud relay.
4. Verify daemon and cloud relay locally.
5. Merge consumers in dependency order.

