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

var url string = "http://localhost:3700/api"

type LoggerApi struct {
	client *http.Client
	id     int64
	token  string
	url    string
}

func NewWithUrl(token string, id int64, url string) *LoggerApi {
	return &LoggerApi{
		client: &http.Client{},
		id:     id,
		token:  token,
		url:    url,
	}
}

func New(token string, id int64) *LoggerApi {
	return NewWithUrl(token, id, url)
}

func (l *LoggerApi) Check() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, l.url+"/check/"+l.token, bytes.NewBuffer([]byte{}))
	if err != nil {
		return false, vanerrors.NewWrap(ErrorCreatingRequest, err, vanerrors.EmptyHandler)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return false, vanerrors.NewWrap(ErrorDoingRequest, err, vanerrors.EmptyHandler)
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return false, vanerrors.NewWrap(ErrorDecodingResponse, err, vanerrors.EmptyHandler)
	}

	if res.Ok || res.Message == "token does not exist" || res.Message == "token exists" {
		return res.Ok, nil
	}

	return false, vanerrors.New(vanerrors.ErrorData{Name: res.Message, Description: res.Description, Code: res.StatusCode}, vanerrors.Options{ShowDescription: true, ShowCode: true}, vanerrors.EmptyLoggerOptions)
}

func (l *LoggerApi) FastCheck() bool {
	res, _ := l.Check()

	return res
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

	req, err := http.NewRequest(http.MethodPost, l.url+"/send", bytes.NewBuffer(body))
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

func (l *LoggerApi) GetUrl() string {
	return l.url
}

func (l *LoggerApi) SetUrl(u string) {
	l.url = u
}

func (l *LoggerApi) SetId(id int64) {
	l.id = id
}

func (l *LoggerApi) GetId() int64 {
	return l.id
}

func (l *LoggerApi) UpdateClient() {
	l.SetClient(&http.Client{})
}

func (l *LoggerApi) SetClient(c *http.Client) {
	l.client = c
}
