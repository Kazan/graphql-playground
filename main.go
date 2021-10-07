package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/kazan/graphql-playground/pkg/codec"
	"github.com/kazan/graphql-playground/schema"
)

func main() {
	s, err := schema.NewSchema()
	if err != nil {
		panic(err)
	}
	codec := codec.NewJSONCodec()

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Handling...")
			req := schema.GraphQLRequest{}

			err := codec.Decode(r.Body, &req)
			if err != nil {
				w.Write([]byte(`Error decoding request`))
			}

			res := graphql.Do(graphql.Params{
				Context:        context.Background(),
				Schema:         s,
				RequestString:  req.Query,
				VariableValues: req.Variables,
				OperationName:  getOperationName(req),
			})

			codec.Encode(w, res)
		},
		// api.NewHandler(schema, graphql.Do, commonApi.NewJSONCodec(), cfg.Origin, mdws...).Handle)
	)

	server := &http.Server{
		Addr:     ":9192",
		Handler:  mux,
		ErrorLog: log.New(os.Stderr, "http: ", log.LstdFlags),
	}

	log.Fatal(server.ListenAndServe())
}

func getOperationName(gr schema.GraphQLRequest) string {
	if gr.OperationName != "" {
		return gr.OperationName
	}

	AST, err := parser.Parse(parser.ParseParams{
		Source: gr.Query,
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   false,
		},
	})
	if err != nil {
		return ""
	}

	for _, definition := range AST.Definitions {
		switch definition := definition.(type) {
		case *ast.OperationDefinition:
			if definition.GetName() != nil {
				return definition.GetName().Value
			}
		default:
			continue
		}
	}

	return ""
}
