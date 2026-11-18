# Goblin

Goblin is a CLI tool for bulk analysis, extraction and assessment of multiple file types within an isolated container. The program ingests small to large volumes of files and utilises container instances to perform the work. All execution takes place within a running Incus container, currently tested using an Alpine image. The environment has limited privileges, only what is needed to perform the process.

The main purpose is found when there is nesting, extensions requiring extra steps to view content, and a bulk amount of files that you do not want extracted on your host PC.

## Architecture

- **Go** — orchestrator. REPL loop, Incus API calls, file ingestion from `./queue`.
- **Rust** — worker binaries (per file type). Pushed into the container, perform analysis, produce a report. *(In development.)*

The Incus client communicates with the Incus daemon (`incusd`) over a Unix socket.

## Commands

| Command | Description |
|---------|-------------|
| `create` | Creates a new instance, e.g. `worker-alpine-1` |
| `delete` | Deletes a specified instance (1..10) |
| `start` | Starts a specified instance (1..10) |
| `stop` | Stops a specified instance (1..10) |
| `push` | Pushes files from `./queue` to a specified instance |
| `pull` | Pulls a file from a specified instance back to the host (requires user confirmation) |
| `list` | Lists all current instances with name, status, type |
| `queue` | Shows files currently in `./queue` |
| `help` | Prints available commands |
| `exit` | Shuts down and exits |

Instances are named `worker-alpine-N` where N is 1..10. Container count, name, and status are fetched from the Incus API at the start of every run.

## File Types (planned)

- zip
- tar
- gz
- xz
- bz2

Each type gets its own Rust executable. One file in, one report out.

## References

- Incus client: <https://pkg.go.dev/github.com/lxc/incus/client>
- Incus API: <https://pkg.go.dev/github.com/lxc/incus/v7/shared/api>
- Incus docs: <https://linuxcontainers.org/incus/docs/main/>   
