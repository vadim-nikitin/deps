deps
====
### Dependency management utility for Linux executables written in Go

To build deps execute these commands in your $GOPATH directory:

```bash
git clone https://github.com/vadim-nikitin/deps.git
go install deps
```

Now you are able to use deps.

### Basic usage

To gather all dependencies for an executable use:

```bash
deps file path
```

where *file* is an executable, *path* is a directory where dependency tree will be created.

### Known issues

* For now deps only can gather dependencies, no any other actions
* The circular dependencies are not yet handled

Enjoy!
