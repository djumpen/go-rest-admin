package api

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

// we need to store request objects to print proper error field on validation failure
// usage in middleware/handler/error_handler.go
var registeredRequests map[string]interface{}

func init() {
	registeredRequests = make(map[string]interface{})

	registerRequests(
		ImageReq{},
		MainSettingsReq{},
	)
}

// Add requests to registry
func registerRequests(r ...interface{}) {
	for _, v := range r {
		name := reflect.TypeOf(v).String()
		registeredRequests[strings.Split(name, ".")[1]] = v
	}
}

// Returns registered requests
func RegisteredRequests() map[string]interface{} {
	return registeredRequests
}

// Prints alert if not all requests was registered via registerRequests()
func PrintRegisteredRequestsReminder() {
	folder := "./api/"
	postfix := "Req"
	requests, err := scanRequests(folder, postfix)
	if err != nil {
		log.Fatal(err)
	}
	match := false
	for _, req := range requests {
		if _, ok := registeredRequests[req]; !ok {
			fmt.Printf("-------- %s not registered --------\n", req)
			match = true
		}
	}
	if match {
		fmt.Println("Register reguest structs in 'api/register_requests.go'")
	}
}

// Scan for all structs with provided postfix
func scanRequests(folder, postfix string) ([]string, error) {
	var requests []string
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, folder+file.Name(), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.TYPE {
					for _, spec := range decl.Specs {
						tspec := spec.(*ast.TypeSpec)
						structName := tspec.Name.String()
						checkLenght := len(structName) - len(postfix)
						if strings.LastIndex(structName, postfix) == checkLenght {
							requests = append(requests, structName)
						}
					}
				}
			}
		}
	}
	return requests, nil
}
