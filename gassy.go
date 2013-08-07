package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "flag"
  "path/filepath"
  "regexp"
  "os"
  sass "github.com/moovweb/gosass"
  "github.com/kylelemons/go-gypsy/yaml"
  "github.com/howeyc/fsnotify"
)

var (
  file = flag.String("f", "config.yml", "YAML config file (defaults to config.yml)")
)

func main() {
  // watchFiles := []string{}
  compileFiles := []string{}
  watchings := map[string]string{}
  watchDirs := []string{}


  flag.Parse()

  // parse the yaml yo

  yml, err := yaml.ReadFile(*file)
  if err != nil {
    log.Fatal("Error reading config file")
  }

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
    s = filepath.Clean(s)
    d = filepath.Clean(d)
    watchDirs = append(watchDirs, s)

    _, compileFiles = getFiles(s)

    for z:=0; z<len(compileFiles); z++ {
      compile(compileFiles[z], d)
      watchings[compileFiles[z]] = d
    }
  }

  params := flag.Args()
  for _, param := range params {
    if param == "build" || param == "b" { // "gassy build"
      fmt.Println("Built. Come again soon.")
    } else if param == "watch" || param == "w" { // default is watch
      watcher, err := fsnotify.NewWatcher()
      if err != nil {
          log.Fatal(err)
      }

      done := make(chan bool)

      go func() {
        for {
          select {
          case ev := <-watcher.Event:
            // log.Println("event:", ev)
            if ev.IsCreate() {
              if finfo, err := os.Stat(ev.Name); err == nil && finfo.IsDir() {
                watcher.Watch(ev.Name)
                log.Println("Added Watch: ", ev.Name)
              }
            } else if ev.IsModify() {
              if finfo, err := os.Stat(ev.Name); err == nil {
                filename := finfo.Name()
                extension := filepath.Ext(filename)
                if (extension == ".scss" || extension == ".sass" || extension == ".css") {
                  // compile everything for now
                  for k, v := range watchings {
                    fmt.Println(k, v)
                    compile(k, v)
                  }
                }
              }
            }
          case err := <-watcher.Error:
            log.Println("error:", err)
          }
        }
        done <- true
      }()

      for i:=0; i<len(watchDirs); i++ {
        err = watchAllDirs(watcher, watchDirs[i])
        if err != nil {
          log.Fatal(err)
        }
      }

      <-done

      watcher.Close()
    }
  }
}




func watchAllDirs(watcher *fsnotify.Watcher, root string) (err error) {
  walkFn := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() {
      watcher.Watch(path)
      log.Println("Added Watch: ", path)
    }
    return nil
  }

  return filepath.Walk(root, walkFn)
}





func getFiles(s string) ([]string, []string) {
  watchFiles := []string{}
  compileFiles := []string{}

  dirList, err := ioutil.ReadDir(s)
  if err != nil {
    log.Fatal("error reading specified directory")
  }

  for x:=0; x<len(dirList); x++ {
    if dirList[x].IsDir() {
      subdirsWatch, subdirsCompile := getFiles(filepath.Join(s, dirList[x].Name()))
      for y:=0; y<len(subdirsWatch); y++ {
        watchFiles = append(watchFiles, subdirsWatch[y])
      }
      for y:=0; y<len(subdirsCompile); y++ {
        compileFiles = append(compileFiles, subdirsCompile[y])
      }
    } else {
      filename := dirList[x].Name()
      extension := filepath.Ext(filename)
      // first find the files we need to compile
      if !(filename[0] == 95) { // ignore files starting with an underscore
        if (extension == ".scss" || extension == ".sass") {
          compileFiles = append(compileFiles, filepath.Join(s, dirList[x].Name()))
        }
      }
      // then find all the files we need to watch for changes
      if (extension == ".scss" || extension == ".sass" || extension == ".css") {
        watchFiles = append(watchFiles, filepath.Join(s, dirList[x].Name()))
      }
    }
  }
  return watchFiles, compileFiles
}




func compile(s string, d string) {
  // compile the sass yo
  ctx := sass.FileContext {
    Options: sass.Options {
      OutputStyle: sass.NESTED_STYLE,
      IncludePaths: make([]string, 0),
    },
    InputPath:    s,
    OutputString: "",
    ErrorStatus:  0,
    ErrorMessage: "",
  }

  // minified version
  ctxMin := sass.FileContext {
    Options: sass.Options {
      OutputStyle: sass.COMPRESSED_STYLE,
      IncludePaths: make([]string, 0),
    },
    InputPath:    s,
    OutputString: "",
    ErrorStatus:  0,
    ErrorMessage: "",
  }

  sass.CompileFile(&ctx)
  sass.CompileFile(&ctxMin)

  if ctx.ErrorStatus != 0 {
    if ctx.ErrorMessage != "" {
      fmt.Print(ctx.ErrorMessage)
    } else {
      fmt.Println("An error occured; no error message available.")
    }
  } else {
    re := regexp.MustCompile("scss|sass")
    name := re.ReplaceAllString(filepath.Base(s), "css")
    nameMin := re.ReplaceAllString(filepath.Base(s), "min.css")
    // write out un-minified file
    err := ioutil.WriteFile(filepath.Join(d, name), []byte(ctx.OutputString), 0644)
    if err != nil {
      panic(err)
    }
    // write out minified file
    err = ioutil.WriteFile(filepath.Join(d, nameMin), []byte(ctxMin.OutputString), 0644)
    if err != nil {
      panic(err)
    }
  }
}