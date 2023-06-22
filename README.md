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

Building your own tools creates an intimate understanding of your computing environment.
I enjoy the comfort and familiarity of being at home on my computer.

## Additional Information

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

The order of the yaml file is not bearing on the order of operations. Casa will always execute in
the following order.

* Create instructions in the order they are specified
* Link instructions in the order they are specified

## TODO

Expect perhaps a couple more items to be added to this list,
but not a lot in order to keep scope small.

- [x] Create top level instruction
- [x] Link top level instruction
- [x] IF property for a link
- [x] Clean top level instruction (removes specified symlinks)
- [ ] Command line arguments (custom yaml file name, custom working directory?)
- [ ] Testing
