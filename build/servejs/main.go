package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var port = flag.Int("port", 8080, "port to serve")

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	addr := "localhost:" + strconv.Itoa(*port)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to getwd: %w", err)
	}
	exampleDir := filepath.Dir(wd)
	err = os.Chdir(exampleDir)
	if err != nil {
		return fmt.Errorf("failed to chdir: %w", err)
	}
	fmt.Printf("serving at %v\n", addr)
	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		switch r.URL.Path {
		case "/wasm_exec.js":
			rw.Header().Add("Content-Type", "application/javascript")
			_, _ = rw.Write(wasmExec)
			return
		case "/favicon.ico":
			rw.WriteHeader(404)
			return
		default:
			// Two requests fallthrough: directories and directory/sample-project_jswasm.wasm
			isWasmRequest := strings.HasSuffix(r.URL.Path, "sample-project_jswasm.wasm")
			if !isWasmRequest {
				r.URL.Path = filepath.Join(r.URL.Path, "sample-project_jswasm.wasm")
			}
			// If no wasm file exists to render, report so
			wasm, err := os.Open(filepath.Join(".", r.URL.Path))
			if err != nil {
				rw.WriteHeader(404)
				fmt.Fprintf(rw, "No wasm file found at %v. Run the following script to generate it:\n\n", r.URL.Path)
				subDir := filepath.Dir(filepath.Join(exampleDir, r.URL.Path))
				fmt.Fprintf(rw, "cd %v\n", subDir)
				_, _ = rw.Write([]byte(`(powershell) $Env:GOOS="js"; $Env:GOARCH="wasm"; go build -o sample-project_jswasm.wasm .`))
				_, _ = rw.Write([]byte("\n"))
				_, _ = rw.Write([]byte("(unix)       GOOS=js GOARCH=wasm go build -o sample-project_jswasm.wasm .\n"))
				fmt.Println("no wasm", err)
				return
			}
			defer func() { _ = wasm.Close() }()

			// serve wasm when requested
			if isWasmRequest {
				rw.Header().Add("Content-Type", "application/wasm")
				_, _ = io.Copy(rw, wasm)
				return
			}

			// otherwise serve the wrapper index that will run the wasm
			t, err := template.New("script").Parse(string(indexHTML))
			if err != nil {
				rw.WriteHeader(404)
				fmt.Println("template error:", err)
				return
			}
			err = t.Execute(rw, struct {
				WasmScript string
			}{
				WasmScript: r.URL.Path,
			})
			if err != nil {
				rw.WriteHeader(404)
				fmt.Println("template exec error:", err)
				return
			}
		}
	}))
	// nolint:gosec // ok for now.
	return http.ListenAndServe(addr, nil)
}

//go:embed index.html.tpl
var indexHTML []byte

//go:embed wasm_exec.js
var wasmExec []byte
