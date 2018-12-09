package shoppinglist

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	host 	string
	port 	int
	url 	string
	client 	http.Client
}

type ShoppingItem struct {
	Item 	string 	`json:"item,omitempty"`
}

func NewClient(host string, port int) *Client {
	url := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	return &Client{host: host, port: port, url: url, client: http.Client{}}
}


func (c Client) makeRequest(items []string) (result []byte, err error) {
	endpoint := strings.Join([]string{c.url, "items"}, "/")

	var jsonItems []ShoppingItem
	for _, item := range items {
		jsonItems = append(jsonItems, ShoppingItem{Item: item})
	}
	b, err := json.Marshal(jsonItems)
	if err != nil {
		log.Print(err)
		return
	}
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(b))
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


func (c Client) AddItems(s []string) (err error) {
	_, err = c.makeRequest(s)
	return
}