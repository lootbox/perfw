package main

import (
	"flag"
	"github.com/lootbox/tripoli"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

type App struct {
	clean bool
	files int
	bytes int
	pproc int
	wkdir string
	templ string
	fbody byte
}

func (o App) workerWriter(amountOfBytes int) {

	f, err := ioutil.TempFile(o.wkdir, o.templ)

	if err != nil {
		log.Fatal(err)
		return
	} else {
		if o.clean {
			defer os.Remove(f.Name())
		}
	}

	for i := 0; i < amountOfBytes; i++ {
		if _, err := f.Write([]byte(".")); err != nil {
			log.Fatal(err)
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (o App) diskWriteTest() {

	data := make([]interface{}, o.pproc)

	for i := range data {
		data[i] = o.bytes
	}

	tripoli.Run(o.workerWriter, o.pproc, data)
}

func New() App {
	var bytes = flag.Int("b", 1024*1024, "Size that will be written in one file in bytes")
	var clean = flag.Bool("c", true, "Recording with cleanup (true) or eat free space")
	var files = flag.Int("f", 9001, "Amount of files that will be created. If '-c=true', created files amount in any time will be constant and will be equal passed to '-p' parameter")
	var pproc = flag.Int("p", runtime.NumCPU()*9+1, "Parallel threads by NumCPU() * 9 + 1")
	var wkdir = flag.String("w", "/tmp", "Working directory: where to write files")
	var templ = flag.String("t", "lootbox******", "Temporary files template")

	flag.Parse()

	return App{
		clean: *clean,
		files: *files,
		bytes: *bytes,
		pproc: *pproc,
		wkdir: *wkdir,
		templ: *templ,
	}
}

func main() {
	app := New()
	app.diskWriteTest()
}
