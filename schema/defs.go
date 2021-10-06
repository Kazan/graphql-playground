package schema

import (
	"fmt"
	"time"

	graphql "github.com/graphql-go/graphql"
)

func NewSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"query1": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "QueryOneOuterField",
						Fields: graphql.Fields{
							"outerField": &graphql.Field{
								Type: graphql.String,
							},
							"embedded": &graphql.Field{
								Type: graphql.NewObject(graphql.ObjectConfig{
									Name: "QueryOneInnerObject",
									Fields: graphql.Fields{
										"innerField1": &graphql.Field{
											Type: graphql.String,
										},
										"innerField2": &graphql.Field{
											Type: graphql.String,
										},
									},
								}),
								Resolve: func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(3 * time.Second)
										fmt.Println("INNER RESOLVING...", i)
									}

									return "InnerResolver", nil
								},
							},
						},
					}),
					Args: graphql.FieldConfigArgument{
						"anything": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						for i := 0; i < 5; i++ {
							time.Sleep(3 * time.Second)
							fmt.Println("ROOT RESOLVING...", i)
						}
						return "RootResolver", nil
					},
				},
			},
		}),
	})
}
