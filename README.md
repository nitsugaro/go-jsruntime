# GO-JSRUNTIME

Basic extension from [goja](https://github.com/dop251/goja) library.

## Features

- require modules
- asynchronous code snippet
- default JSON storage

```bash
go get github.com/nitsugaro/go-jsruntime@latest
```

## Usage

```go
import jsrun "github.com/nitsugaro/go-jsruntime"
import "github.com/nitsugaro/go-utils/encoding"

var logger = func(msg string) {
	fmt.Println("JS Log:", msg)
}

var testHTTPRequest = func() {
	resp, _ := http.Get("https://httpbin.org/get")
	defer resp.Body.Close()
	fmt.Println("HTTP GET:", resp.StatusCode)
}

var asyncScript = &jsrun.Script{
	Name: "http-async",
	Type: "async",
	CodeBase64: encoding.EncodeBase64([]byte(`
		logger("Async request triggered")
		http()
	`)),
}

var lib1 = &jsrun.Script{
	Name: "library-1",
	Type: "library",
	CodeBase64: encoding.EncodeBase64([]byte(`
		let count = 0
		exports.inc = () => count++
		exports.get = () => count
	`)),
}


/* ########### MAIN ############# */

main := `
	const lib = require("library-1")

	for (let i = 0; i < 5; i++) lib.inc()

	logger("Count is: " + lib.get())
	executeAsync("http-async", { logger, http })
`

manager, store := jsrun.NewDefaultStorage("scripts")
store.Save(lib1)
store.Save(asyncScript)

compiled, _ := manager.CompileScript("main", main)

manager.ExecuteWithBindings(compiled, map[string]interface{}{
	"logger": logger,
	"http":   testHTTPRequest,
})
```