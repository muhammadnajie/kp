package resources

type Success struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type Failed struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
