package schema

// GraphQLRequest defines a GraphQL request
type GraphQLRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

// GraphQLResponse defines a GraphQL response
type GraphQLResponse struct {
	Data   interface{}            `json:"data"`
	Errors []GraphQLResponseError `json:"errors"`
}

// GraphQLResponseError defines a GraphQL response error
type GraphQLResponseError struct {
	Message    string
	Extensions struct {
		Code       string
		Stack      []string
		Violations []GraphQLViolation
	}
}

// GraphQLViolation defines a GraphQL violation in response
type GraphQLViolation struct {
	PropertyPath string
	Message      string
}
