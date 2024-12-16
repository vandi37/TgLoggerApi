package tgloggerapi

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vandi37/vanerrors"
)

const (
	ErrorMarshaling       = "error marshaling"
	ErrorCreatingRequest  = "error creating request"
	ErrorDoingRequest     = "error doing request"
	ErrorDecodingResponse = "error decoding response"
)

var link string = "http://localhost:3700/send"

type LoggerApi struct {
	client *http.Client
	id     int64
	token  string
}

func New(token string, id int64) *LoggerApi {
	return &LoggerApi{
		client: http.DefaultClient,
		id:     id,
		token:  token,
	}
}

func (l *LoggerApi) GetId() int64 {
	return l.id
}

func (l *LoggerApi) Write(p []byte) (n int, err error) {
	data := Request{
		Token: l.token,
		Id:    l.id,
		Text:  string(p),
	}
	body, err := json.Marshal(data)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorMarshaling, err, vanerrors.EmptyHandler)
	}

	req, err := http.NewRequest(http.MethodPost, link, bytes.NewBuffer(body))
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorCreatingRequest, err, vanerrors.EmptyHandler)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorDoingRequest, err, vanerrors.EmptyHandler)
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorDecodingResponse, err, vanerrors.EmptyHandler)
	}

	if res.Ok {
		return len(p), nil
	}

	return 0, vanerrors.New(vanerrors.ErrorData{Name: res.Message, Description: res.Description, Code: res.StatusCode}, vanerrors.Options{ShowDescription: true, ShowCode: true}, vanerrors.EmptyLoggerOptions)
}
