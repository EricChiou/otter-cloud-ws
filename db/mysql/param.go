package mysql

// SQLParams sql parameters
type SQLParams struct {
	kv map[string]string
}

// Add value with key
func (params *SQLParams) Add(key string, value string) {
	params.kv[key] = value
}

// Get value with key
func (params *SQLParams) Get(key string) string {
	return params.kv[key]
}

// Remove value which key equal to input parameter
func (params *SQLParams) Remove(key string) {
	delete(params.kv, key)
}

// SQLParamsInstance get sqlParams instance
func SQLParamsInstance() SQLParams {
	return SQLParams{kv: map[string]string{}}
}
