package test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	jsrun "github.com/nitsugaro/go-jsruntime"
	"github.com/nitsugaro/go-utils/encoding"
)

func Logger(message string) {
	fmt.Println(message)
}

func testHTTPRequest() {
	resp, err := http.Get("https://httpbin.org/delay/5")
	if err != nil {
		log.Println("Error haciendo GET:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("status code:", resp.StatusCode)
}

var asyncFunc = &jsrun.Script{
	Name: "http-async",
	Type: "async",
	CodeBase64: encoding.EncodeBase64([]byte(`
		if(typeof async !== "undefined") {
			logger("calling async again!")
			async("http-async", { logger, http })
		}

		person.hi()

		logger("executing request")
		http()
		logger("end request")
	`)),
}

var library1 = &jsrun.Script{
	Name: "library-1",
	Type: "library",
	CodeBase64: encoding.EncodeBase64([]byte(`
		logger("init Library-1")

		let localVariable = 1
		exports.getLocalVariable = () => {
			return localVariable
		}

		exports.incrementLocalVariable = () => {
			localVariable += 1
		}
	`)),
}

var library2 = &jsrun.Script{
	Name: "library-2",
	Type: "library",
	CodeBase64: encoding.EncodeBase64([]byte(`
		logger("init Library-2")
		const library1 = require("library-1")

		library1.incrementLocalVariable()

		class Person {
			hi() {
				logger("Hi, I'm a Person")
			}
		}

		module.exports = {
			Person
		}
	`)),
}

var mainScript = `
		logger("executing main script")
		const library1 = require("library-1")
		const { Person } = require("library-2")

		logger("local varible: " + library1.getLocalVariable())
		if(library1.getLocalVariable() != 2) {
			logger("expected localVariable to be equal to 2 and got " + library1.getLocalVariable())
			throw new Error()
		}

		const limit = 10
		function loop() {
			for(let i = 0; i < limit; i++) logger(i)
		}

		const myPerson = new Person()

		executeAsync("http-async", { logger, http, async: executeAsync, person: myPerson })
		executeAsync("http-async", { logger, http, person:myPerson })

		callbacks.Add(1)
		
		loop()
		myPerson.hi()
		library1.getLocalVariable()
	`

func TestMain(t *testing.T) {
	var scriptManager, scriptStorage = jsrun.NewDefaultStorage("scripts")

	scriptStorage.Save(asyncFunc)
	scriptStorage.Save(library1)
	scriptStorage.Save(library2)

	prog, _ := scriptManager.CompileScript("main", mainScript)
	val, err := scriptManager.ExecuteWithBindings(prog, map[string]interface{}{"logger": Logger, "http": testHTTPRequest})
	if err != nil {
		t.Error(err.Error())
	}

	if val.ToNumber().ToInteger() != 2 {
		t.Errorf("val must be equal to 2")
	}

	time.Sleep(10 * time.Second)

	os.RemoveAll("scripts")
}
