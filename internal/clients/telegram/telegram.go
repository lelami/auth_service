package telegram

import (
	e "authservice/internal/helpers"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}
func newBasePath(token string) string {
	return "bot" + token
}
func (c *Client) SendCode(code string) error {
	return nil
}
func (c *Client) Updates(offset, limit int) ([]Update, error) {

	q := &url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}
func (c *Client) SendMessage(chatId int, text string) error {
	q := &url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.WrapIfErr("can't send message", err)
	}
	return nil
}
func (c *Client) doRequest(method string, query *url.Values) (data []byte, err error) {
	//defer func() { err = e.WrapIfErr("can't do request:", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = query.Encode()
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Println("error closing response body", cerr)
		}
	}()
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
