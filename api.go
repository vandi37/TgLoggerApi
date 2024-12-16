package tgloggerapi

// A response
type Response struct {
	Ok          bool   `json:"ok"`
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// A request
type Request struct {
	Token string `json:"token"`
	Id    int64  `json:"id"`
	Text  string `json:"text"`
}
