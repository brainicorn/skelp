## skelp bashme

Creates a bash completion file for skelp

### Synopsis


Creates a bash completion file for skelp.
By default the completion file is written to /etc/bash_completion.d/ using sudo.

```
skelp bashme [flags]
```

### Options

```
  -h, --help            help for bashme
      --no-sudo         will try to write the completion file without using sudo
      --output string   path to the directory where the completion file will be written (default "/etc/bash_completion.d/")
```

### Options inherited from parent commands

```
      --homedir string    path to override user's home directory where skelp stores data
      --no-color          turn off terminal colors
      --quiet             run in 'quiet mode'
      --skelpdir string   override name of skelp folder within the user's home directory
```

### SEE ALSO
* [skelp](skelp.md)	 - A commandline tool for generating skeleton projects

