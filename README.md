# Gassy
### *go with sass*

Gassy is a command line tool used to watch directories containing Sass files and automatically compile them to CSS whenever a file change is detected.

### BUILDING

First install Go (http://golang.org/doc/install) if you don't have it already.

Next, install libsass (https://github.com/hcatlin/libsass)

After that, `go get` the following three libraries:

* github.com/moovweb/gosass
* github.com/kylelemons/go-gypsy/yaml
* github.com/howeyc/fsnotify

Finally, run a `go install gassy.go` and you're all set


### CONFIG

To get started, first change the config.yml file to list the directories you would like to watch and the directories where you would like to build the CSS.

It should be in this format:

```YAML
config:
  - source:  test/source1
    dest: test/compiled
  - source:  test/source2
    dest: test/compiled
```

You can have as many watched directories as you like (*the sky is the limit*).

### RUNNING GASSY
After compiling, running gassy is as simple as

    gassy -f config.yml watch

To just build your stylesheets but not watch, run:

    gassy -f config.yml build

If you use the default config.yml (and it's in the same directory as the executable), you can also ignore the -f param and simply run:

    gassy watch