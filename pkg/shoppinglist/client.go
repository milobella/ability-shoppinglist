package shoppinglist

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

type Client struct {
    host    string
    port    int
    url     string
    client  http.Client
}

type ShoppingItem struct {
    Item    string  `json:"item,omitempty"`
}

func NewClient(host string, port int) *Client {
    url := fmt.Sprintf("http://%s:%d", host, port)
    return &Client{host: host, port: port, url: url, client: http.Client{}}
}

func (c Client) makeRequest(method string, endpoint string, input []byte) (result []byte, err error) {

    req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(input))
    if err != nil {
        log.Print(err)
        return
    }

    resp, err := c.client.Do(req)
    if err != nil {
        log.Print(err)
        return
    }
    output, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    if err != nil {
        log.Print(err)
        return
    }

    return output, err
}


func (c Client) AddItems(items []string) (err error) {
    endpoint := strings.Join([]string{c.url, "items"}, "/")

    // Build items and serializes it
    var jsonItems []ShoppingItem
    for _, item := range items {
        jsonItems = append(jsonItems, ShoppingItem{Item: item})
    }
    b, err := json.Marshal(jsonItems)
    if err != nil {
        log.Print(err)
        return
    }

    // Proceed the request
    _, err = c.makeRequest("PUT", endpoint, b)
    return
}

func (c Client) GetItems() (result []string, err error) {
    endpoint := strings.Join([]string{c.url, "items"}, "/")

    // Proceed the request
    resp, err := c.makeRequest("GET", endpoint, nil)

    // Deserialize response and build items
    var items []ShoppingItem
    err = json.Unmarshal(resp, &items)

    if err != nil {
        log.Print(err)
        return
    }

    // Get items' denomination
    for _, element := range items {
        result = append(result, element.Item)
    }

    return
}