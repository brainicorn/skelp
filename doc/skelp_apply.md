## skelp apply

Apply a template to the current directory

### Synopsis


Apply a template to the current directory

```
skelp apply [git-url|file-path|alias] [flags]
```

### Options

```
  -d, --data string     path to a json data file for filling in template data
      --dry-run         just gather data, no generation (for testing)
  -f, --force           force overwriting of files without asking
  -h, --help            help for apply
      --offline         turns off auto-downloading/updating of templates
  -o, --output string   path to the directory where the template should be applied (default "current directory")
```

### Options inherited from parent commands

```
      --homedir string    path to override user's home directory where skelp stores data
      --no-color          turn off terminal colors
  -q, --quiet             run in 'quiet mode'
      --skelpdir string   override name of skelp folder within the user's home directory
```

### SEE ALSO
* [skelp](skelp.md)	 - A commandline tool for generating skeleton projects

