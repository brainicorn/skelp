[![Build Status](https://travis-ci.org/brainicorn/skelp.svg?branch=master)](https://travis-ci.org/brainicorn/skelp)
[![codecov](https://codecov.io/gh/brainicorn/skelp/branch/master/graph/badge.svg)](https://codecov.io/gh/brainicorn/skelp)
[![Go Report Card](https://goreportcard.com/badge/github.com/brainicorn/skelp)](https://goreportcard.com/report/github.com/brainicorn/skelp)
[![GoDoc](https://godoc.org/github.com/brainicorn/skelp?status.svg)](https://godoc.org/github.com/brainicorn/skelp)

![skelp logo](https://rawgit.com/brainicorn/skelp/master/doc/skelp-logo.svg)

# skelp #

#### skeleton project generator and template runner

## Overview


Skelp is both a command line tool and a golang library for generating project files. It can be used to generate entire projects from scratch as well as adding files to existing projects by applying smaller template sets.

Skelp can use templates that are hosted in git-compatible source code repositories as well as templates stored on your local file system.

Although skelp itself is written in Go, it can be used to generate any type of files in any language/format. On top of that, the skelp tool itself is a stand-alone binary that does not require Go to be installed to run. Simply download the proper version for your platform and run it.

## Installation

### Using Go
```bash
go get -u github.com/brainicorn/skelp
```

### Binary install without Go

Stand-alone binary releases can be downloaded from the [releases page](https://github.com/brainicorn/skelp/releases)

Simply download the correct binary for your operating system and run it. Check your platform specific documentation for how to get the binary into your path.

### Quick start for the really impatient
If you just want to give it a go (pun-intended), after installation you can check that skelp runs properly by typing
```bash
skelp version
```

If you want to test out generating something, navigate to a temporary or new folder and type:
```bash
skelp apply https://github.com/brainicorn/skelp-simple-readme
```
This will download a simple template that generates a very basic README file.

And you can always type ```skelp help ``` to see all of skelp's commands/options.

#### For detailed information:
- [User Guide](doc/user-guide.md)
- [Template Guide](doc/template-guide.md)
- [Library Guide](doc/library-guide.md)

## Features

### For Users
- Easy subcommand-based CLI
- Easy to use terminal prompts - keyed input, select boxes, multiple selections, etc.
- Download/Update and apply templates from a repo in a single command
- Local templates
- Auto-updating of templates
- Local template caching
- Supply reusable data files as full/partial input (allows quickly skipping prompts)
- Full control over overwrites
- Offline Mode
- Public and Private repository support
- SSH and Basic repository authentication
- Template aliasing - no need to remember long URLs
- Alias management
- Bash Completion

### For Template Authors
- Full [golang template](https://golang.org/pkg/text/template/) support
- Variables **everywhere**: within templates, file names, folder names, default values, variable names...
- Built-in [golang functions](https://golang.org/pkg/text/template/#hdr-Functions) support
- Full [sprig functions](https://github.com/Masterminds/sprig) support
- JSON-based project descriptor
  - json-schema is provided
  - validation tools are provided
- Ability to define variables/input as:
  - Simple keyed input
  - Single-select boxes
  - Multi-select boxes
  - Multi-value input (e.g. 'add another?')
  - Required values
  - Masked input (passwords)
- Ability to customize prompts per-variable
- Define min/max length/values for input
- Ordered input gathering for using previously inputted values in variables/defaults
- Ability to use data files for "promptless" template testing

### For Golang Developers
- Anything the CLI can do can be done programatically
  - Automate any/everything skelp can do
  - Create a custom CLI if you wish
- Ability to customize key filenames / directories (home, skelp.json, etc)
- Can be used as a lower-level library for simply applying templates across a directory tree
  - doesn't require an skelp-specific concepts/files/directories
  - basically like executing a single template in go, only traversing an entire directory of templates
- "Pluggable" customization of the way input data is gathered (prompts, files, basically anything)
- "Pluggable" customization of the way overwriting is determined

## Basic Use

This section just gives a basic overview of the most commonly used commands.

For full command details and options, see [The User Guide](doc/user-guide.md)

### Enabling Bash Completion

Skelp ships with a command that can add bash completion support for itself.

This is useful if you want to use <tab> to see the available commands/options while using skelp.

To enable bash completion, simply run:

``` skelp bashme ```

**NOTE:** this command uses sudo to copy the completion file to the proper location so it gets picked up by your terminal and may ask you for your sudo password.

After running this command, you'll need to restart your terminal for the changes to take affect.

### Applying Templates

Skelp can apply a template to an empty directory or an existing directory with existing files. The latter is useful for adding smaller pieces to an existing project, like say, a README file.

In either case, you should navigate **into** the root of the directory where you'd like to apply the template as generation output defaults to the current working directory.

#### From a Repository

To apply a template hosted in a repository, simply use:

``` skelp apply <template rpo url> ```

for example:

``` skelp apply https://github.com/brainicorn/skelp-simple-readme ```

The repository url can be any valid git url including https and ssh urls.

#### From a Local Directory

Skelp can also be use templates on your local computer by simply pointing it at the directory holding the template.

Again, change into the directory where you'd like the template applied and type:

``` skelp apply <some template dir> ```

for example:

``` skelp apply ~/templates/my-custom-template ```

### Creating Aliases

If you get tired of always copying/pasting a long repository url when applying templates, you can create a shortcut or "alias" to the template url.

To do so, run:

``` skelp alias add <alias> <template url>```

for example:

``` skelp alias add readme https://github.com/brainicorn/skelp-simple-readme ```

Once an alias has been registered, you can use it in place of the template url when running the apply command.

``` skelp apply readme ```

You can see the current aliases registered by running:

``` skelp alias list ```

## Templates, Templates, Templates

Looking for a template? Have an awesome template you'd like to contribute?

Skelp is a brand new project and we're trying to gather quality templates as fast as we can...

If you're looking for a template to use with skelp, check out our [List of Known Templates](doc/template-list.md)

If you have a template that you'd like to contribute and have listed, please [create an issue](https://github.com/brainicorn/skelp/issues) and include the following:
- URL of the template repository
- Name of the template
- Short description of the template
- (optional) URL to an avatar image for your template
- (optional) Any notes you feel users should know about using your template

## Contributing to Skelp
Pull requests, issues and comments welcome. Please read our [Code of Conduct](doc/CODE_OF_CONDUCT.md)

1. Fork it
1. Create your feature branch (git checkout -b my-new-feature)
1. Commit your changes (git commit -am 'Add some feature')
1. Push to the branch (git push origin my-new-feature)
1. Create new Pull Request **targeting the develop branch**

Since skelp is a command-line tool that can be installed via ```go get``` we do stage all of our development in a branch named "develop" before releasing/merging to the master branch. This ensures that the master branch always represents the latest released version and new development doesn't break current users.

Given this, when making pull requests:
- Ensure you have tests for new features and bug fixes
- Separate unrelated changes into multiple pull requests
- All pull requests should target the **develop** branch

If you're looking for ideas, please look through [the existing issues](https://github.com/brainicorn/skelp/issues)

If you're making a small new fix, create and issue and reference it in your pull request.

If you're making a larger change, make sure you first start a discussion by creating an issue and explaining the intended changes.

## License

Skelp is released under the Apache 2.0 license. See [LICENSE.txt](LICENSE.txt)
