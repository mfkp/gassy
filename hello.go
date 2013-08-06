package main

import (
  // "os"
  "fmt"
  "log"
  // "io/ioutil"
  // sass "github.com/moovweb/gosass"
  "github.com/kylelemons/go-gypsy/yaml"
  "flag"
)

type T struct {
    A string
    B struct{C int; D []int ",flow"}
}

type Lister interface {
   List() yaml.Node
}

var (
  file = flag.String("file", "config.yml", "(Simple) YAML file to read")
)

func main() {
  // fmt.Printf("Hello, world.\n")
  // args := os.Args

  // if len(args) < 2 {
  //   fmt.Println("Usage: gass [INPUT FILE]")
  //   os.Exit(1)
  // }

  // ctx := sass.FileContext {
  //   Options: sass.Options {
  //     OutputStyle: sass.NESTED_STYLE,
  //     IncludePaths: make([]string, 0),
  //   },
  //   InputPath:    args[1],
  //   OutputString: "",
  //   ErrorStatus:  0,
  //   ErrorMessage: "",
  // }

  // sass.CompileFile(&ctx)

  // if ctx.ErrorStatus != 0 {
  //   if ctx.ErrorMessage != "" {
  //     fmt.Print(ctx.ErrorMessage)
  //   } else {
  //     fmt.Println("An error occured; no error message available.")
  //   }
  // } else {
  //   fmt.Print(ctx.OutputString)
  // }

//go-gypsy

  flag.Parse()



  params := flag.Args()
  for _, param := range params {
    // parse the yaml yo
    yml, err := yaml.ReadFile(*file)
    if err != nil {
      log.Fatalf("readfile(%q): %s", *file, err)
    }
  }
  params[0]

  count, err := yml.Count("config")

  for i:=0; i<count; i++ {
    s, err := yml.Get(fmt.Sprintf("config[%d].source", i))
    if err != nil {
      panic(err)
    }
    d, err := yml.Get(fmt.Sprintf("config[%d].dest", i))
    if err != nil {
      panic(err)
    }
    fmt.Println("source: ", s)
    fmt.Println("dest: ", d)
  }

}