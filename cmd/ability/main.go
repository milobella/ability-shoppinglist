package main

import (
    "gitlab.milobella.com/milobella/abilities/ability-sdk-go/pkg/ability"
)

// fun main()
func main() {
    server := ability.NewServer(10200)
    server.RegisterIntent("hello_world", helloWorldHandler)
    server.Serve()
}

func helloWorldHandler(req ability.Request, resp *ability.Response) {
    resp.Nlg.Sentence = "Hello {{world}}"
    resp.Nlg.Params = map[string]string{"world": "world"}
}
