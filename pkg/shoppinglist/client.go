package shoppinglist

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	host string
	port int
	url string
	client http.Client
}

func NewClient(host string, port int) *Client {
	url := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	return &Client{host: host, port: port, url: url, client: http.Client{}}
}


func (c Client) makeRequest(item string) (result []byte, err error) {
	endpoint := strings.Join([]string{c.url, "items"}, "/")

	str := fmt.Sprintf("[{\"item\": \"%s\"}]", item)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(str)))
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	return body, err
}


func (c Client) AddItem(s string) (err error) {
	_, err = c.makeRequest(s)
	return
}