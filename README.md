# CASA

Very much **WIP**. Not intended for use yet.

> Casa is the Spanish word for home.
> It can evoke a sense of warmth, comfort, and familiarity.
> Invoking the idea of managing and organizing one's digital home or personal environment.

A simple dotfile symlink program inspired by [dotbot](https://github.com/anishathalye/dotbot) with a small set of features.

## Motivation

### Simple Scope

I used dotbot prior to link dotfiles. However i used a small subset of dotbot features.
Casa takes just the parts i use, and makes it easier and lighter to get things running.

### Binary distribution

I like the way dotbot works but on some machines installing dotbot created unwanted overhead.
Either because i didn't want or need python installed. Or complications with venv.

### DIY Philosophy

Building your own tools creates an intimate understanding
and familiarity of your computing environment.

## Additional Information

Casa requires a config file to operate, you can pass this in with `--config` or `-c`.

A configuration has 3 root instructions.

* Create - create directories.
* Link - symlink files.
* Clean - Check for dead symlink.

The order of the yaml file (Create, Link, Clean) is not bearing on the order of operations.
Casa will always execute in the following order.

* Clean instructions in the order they are specified
* Create instructions in the order they are specified
* Link instructions in the order they are specified

Here is a sample configuration `install.conf.yaml`.

```yaml
create:
  - ~/.config/abc
  - /tmp/test/test/abc
link:
  - destination: ~/.local/scripts
    path: .local/scripts/test
    if: >
      [ "$(getent passwd $(whoami) | awk -F: '{print $7}')" = "/usr/bin/zsh" ]
  - destination: /tmp/test/test/abc
    path: abc
clean:
  - /tmp/test/test
```

## Root Instructions

### Create

Paths specified under this instruction will be created if they do not exist, if they already exist
they will be skipped.

### Link

Link describes how to symlink a file. The **destination** refers to your symlink,
**path** refers to the file that **destination** will point to.

Links also contain an optional **if** property. This will spawn a shell using `sh`
and run your command, a non zero exit code will result in the link being skipped.

```none
                         symlink will point here
                         this is typically a file inside your dotfiles repo
                        +-----------------------------+
Path------------------->|/home/roland/dotfiles/scripts|
                        +-----------------------------+

                        symlink will be placed here
                        +---------------------------+
Destination------------>|/home/roland/.local/scripts|
                        +---------------------------+
```

### Clean

Clean removes dead symlinks. Some simple rules apply to this.

The working directory is the current directory you are executing the script from.
Not where the script actually might be.

Please keep the working directory in mind,
if you are executing from the home directory then any symlink pointing to
a file inside the home directory could potentially be removed.

```none
+---------------+ for each link+-----------------+
|walk the       |------------->|Is the link dead?|
|clean directory|              +-----------------+
+---------------+                  |       |
                                   v       v
                                 +---+   +---+
                                 |yes|   |no |----+
                                 +---+   +---+    |
                                   |              v
                                   v         +---------+
                      +-----------------+    |no action|
                      |does the link    |    +---------+
                      |point to a file  |         ^
                      |inside the       |         |
                      |working directory|         |
                      +-----------------+         |
                            |       |             |
                            v       v             |
                          +---+   +---+           |
                          |yes|   |no |-----------+
                          +---+   +---+
                            |
                            v
                    +------------------+
                    |delete the symlink|
                    +------------------+
```

## Releasing

Releases can be made from the command line.

```bash
gh release upload v0.0.1 ./casa
gh release create v0.0.1 --repo RolandWarburton/casa --target $(git rev-parse HEAD) --title "Release v0.0.1" --notes "Initial release" 
gh release upload v0.0.1 ./casa
```

## TODO

Expect perhaps a couple more items to be added to this list,
but not a lot in order to keep scope small.

- [x] Create top level instruction
- [x] Link top level instruction
- [x] IF property for a link
- [x] Clean top level instruction (removes specified symlinks)
- [x] Command line arguments (custom yaml file name, custom working directory?)
- [ ] Testing
