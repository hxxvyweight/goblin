
## Concept

Goblin is currently a concept I have for bulk analysis, extraction and assessment of multiple file types, within a isolated container. The program should be able to ingest small to large volumes of files, work concurrently, utilise container instances to perform the work in a useful and optimal manner.
The execution will all take place within a running Incus Container, currently the testing is done utilising an alpine image.
The environment should have limited privileges, only what is needed to perform said process or execution.

Each containers memory limits and usage should be calculated and set before creation, specifically based on the type of or intensity of the task.
On the client side, the program should receive multiple files to ingest based on their path, workers should be created to each receive 1 file and produce one report at a time.

I believe the main purpose of this program is found when there is nesting, extensions requiring extra steps to view content, and a bulk amount of files that you do not wan't extracted on your host PC.

This program will utilise Incus to its required extent depending on needed functionalities, the client documentation is found at <https://pkg.go.dev/github.com/lxc/incus/client>
The documentation for the api is found at <https://pkg.go.dev/github.com/lxc/incus/v7/shared/api> & <https://linuxcontainers.org/incus/docs/main/>

Commands used through REPL loop, making calls to the incus API and server instance *incus.InstanceServer. It is currently possible to use all standard functionalities including:
- **create**:   creates a  new instance e.g [worker-image-1, worker-alpine-2, worker-ubuntu-3]
- **delete**:   deletes a specified instance within range 1..10
- **start**:    starts a specified instance within range 1..10
- **stop**:     stops a specified instance within range 1..10
- **push**:     pushes a file to specified instance within range 1..10
- **pull**:     pulls a file and/or files from specified instance within range 1..10

Incus client communicates with the Incus daemon (incusd) primarily over a Unix Socket.

All file and directory handling will be written in Rust, for each analysis a Rust executable will be 'pushed' into the container along with the specified files. The Rust executable will be utilised through the Incus API client, performing typical linux commands inside the container without having to handle the Rust program on your host. The Golang portion operates as the orchestrator of multiple 'worker-alpine-N' where N is each number specifying an instance. The rust program will do analysis work, reporting and exporting within the container.

An accepted prompt is required by the user to pull files back to the host after using the command 'pull'.

I believe there is strong variety here with use of the incus daemon through your standard CLI, its still possible for fine tuning whether the program is active or stopped.
Container count, name, status is fetched at the start of every run through the incus api, ensuring that count is never lost. This makes it possible for easy execution of any instances, individually.

There also needs to be a seperate Rust executable for each file type listed below (currently starting with these types only):

**File Types**:
- zip
- tar
- gz
- xz
- bz2

