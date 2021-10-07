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
										fmt.Println("Q1 OUTER FIELD RESOLVING...", i)
									}

									return "Q1 OuterFieldResolver", nil
								}).ResolvePromise,
							},
							"anotherOuterField": &graphql.Field{
								Type: graphql.String,
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("Q1 ANOTHER OUTER FIELD RESOLVING...", i)
									}

									return "Q1 AnotherOuterFieldResolver", nil
								}).ResolvePromise,
							},
							"embedded": &graphql.Field{
								Type: graphql.NewObject(graphql.ObjectConfig{
									Name: "QueryTwoInnerObject",
									Fields: graphql.Fields{
										"innerField1": &graphql.Field{
											Type: graphql.String,
										},
										"innerField2": &graphql.Field{
											Type: graphql.String,
											Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
												for i := 0; i < 5; i++ {
													time.Sleep(2 * time.Second)
													fmt.Println("Q1 EMBEDDED FIELD RESOLVING...", i)
												}

												return "Q1 EmbeddedFieldResolver", nil
											}).ResolvePromise,
										},
									},
								}),
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("Q1 INNER RESOLVING...", i)
									}

									return "Q1 InnerResolver", nil
								}).ResolvePromise,
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
				"query2": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "QueryTwoOuterField",
						Fields: graphql.Fields{
							"outerField": &graphql.Field{
								Type: graphql.String,
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("Q2 OUTER FIELD RESOLVING...", i)
									}

									return "Q2OuterFieldResolver", nil
								}).ResolvePromise,
							},
							"anotherOuterField": &graphql.Field{
								Type: graphql.String,
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("Q2 ANOTHER OUTER FIELD RESOLVING...", i)
									}

									return "Q2 AnotherOuterFieldResolver", nil
								}).ResolvePromise,
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
													fmt.Println("Q2 EMBEDDED FIELD RESOLVING...", i)
												}

												return "Q2 OuterFieldResolver", nil
											}).ResolvePromise,
										},
									},
								}),
								Resolve: NewResolverPromise(func(p graphql.ResolveParams) (interface{}, error) {
									for i := 0; i < 5; i++ {
										time.Sleep(2 * time.Second)
										fmt.Println("Q2 INNER RESOLVING...", i)
									}

									return "Q2 InnerResolver", nil
								}).ResolvePromise,
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
							fmt.Println("Q2 ROOT RESOLVING...", i)
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
