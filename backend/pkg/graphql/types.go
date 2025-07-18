package graphql

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Items   interface{} `json:"items"`
	NextKey *string     `json:"nextKey,omitempty"`
}
