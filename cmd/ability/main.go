package main

import (
    "gitlab.milobella.com/milobella/abilities/ability-sdk-go/pkg/ability"
    "gitlab.milobella.com/milobella/abilities/shoppinglist-ability/pkg/shoppinglist"
    "gitlab.milobella.com/milobella/oratio/pkg/anima"
)

var shoppingListClient = shoppinglist.NewClient("http://0.0.0.0", 4848)

// fun main()
func main() {
    server := ability.NewServer(10400)
    server.RegisterIntent("ADD_TO_LIST", addToListHandler)
    server.Serve()
}

func addToListHandler(req ability.Request, resp *ability.Response) {
    // Retrieve only shopping items from NLU entities
    var items []string
    for _, ent := range req.Nlu.Entities {
        if ent.Label == "SHOPITEM" {
            items = append(items, ent.Text)
        }
    }

    // Add these items into the shopping list
    if err := shoppingListClient.AddItems(items); err != nil {
        resp.Nlg.Sentence = "Error adding item to your shopping list"
        return
    }

    // Build the NLG answer
    resp.Nlg.Sentence = "I added {{items}} to your shopping list"
    resp.Nlg.Params = []anima.NLGParam{{
        Name: "items",
        Value: items,
        Type: "enumerated_list",
    }}
}
