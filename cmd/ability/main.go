package main

import (
    "gitlab.milobella.com/milobella/abilities/ability-sdk-go/pkg/ability"
    "gitlab.milobella.com/milobella/abilities/shoppinglist-ability/pkg/shoppinglist"
)

var shoppinglistClient = shoppinglist.NewClient("http://0.0.0.0", 4848)

// fun main()
func main() {
    server := ability.NewServer(10400)
    server.RegisterIntent("ADD_TO_LIST", addToListHandler)
    server.Serve()
}

func addToListHandler(req ability.Request, resp *ability.Response) {
    // TODO: get it from the NLU
    if err := shoppinglistClient.AddItem("des tomates"); err != nil {
        resp.Nlg.Sentence = "Error adding item to your shopping list"
        return
    }

    resp.Nlg.Sentence = "I added {{item}} to your shopping list"
    resp.Nlg.Params = map[string]string{"item": "des tomates"}
}
