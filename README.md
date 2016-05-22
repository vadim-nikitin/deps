deps
====
### Dependency management utility for Linux executables written in Go

It gather dependencies into a directory, saving original dependency directory tree,
so the program can be started from inside a Docker container or chroot environment.
The circular dependencies are handled as well.

To build deps execute these commands:

```bash
git clone https://github.com/vadim-nikitin/deps.git
make
```

Now you are able to use deps.

### Basic usage

To gather all dependencies for an executable use:

```bash
deps EXECUTABLE PATH
```

where *PATH* is a directory where dependency tree will be created.

### Known issues

* For now deps only can gather dependencies, no any other actions

Enjoy!
