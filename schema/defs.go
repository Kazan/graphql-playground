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
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("OUTER FIELD RESOLVING...", i)
									}

									return "OuterFieldResolver", nil
								}).Resolve,
							},
							"anotherOuterField": &graphql.Field{
								Type: graphql.String,
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("ANOTHER OUTER FIELD RESOLVING...", i)
									}

									return "AnotherOuterFieldResolver", nil
								}).Resolve,
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
											Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
												for i := 0; i < 5; i++ {
													time.Sleep(2 * time.Second)
													fmt.Println("EMBEDDED FIELD RESOLVING...", i)
												}

												return "OuterFieldResolver", nil
											}).Resolve,
										},
									},
								}),
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("INNER RESOLVING...", i)
									}

									return "InnerResolver", nil
								}).Resolve,
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
							time.Sleep(2 * time.Second)
							fmt.Println("ROOT RESOLVING...", i)
						}
						return "RootResolver", nil
					},
				},
			},
		}),
	})
}

type fnResolver func(p graphql.ResolveParams) (interface{}, error)

func NewResolverPromise(fn fnResolver) *resolver {
	return &resolver{
		fn: fn,
	}
}

type resolver struct {
	fn fnResolver
}

func (r *resolver) ResolvePromise(p graphql.ResolveParams) (interface{}, error) {
	type result struct {
		data interface{}
		err  error
	}

	ch := make(chan *result, 1)

	go func() {
		defer close(ch)
		res, err := r.fn(p)
		ch <- &result{data: res, err: err}
	}()

	return func() (interface{}, error) {
		res := <-ch

		if res.err != nil {
			return nil, res.err
		}

		return res.data, nil
	}, nil
}

func (r *resolver) Resolve(p graphql.ResolveParams) (interface{}, error) {
	return r.fn(p)
}
