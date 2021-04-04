package main

import (
	"github.com/milobella/ability-sdk-go/pkg/ability"
	"github.com/milobella/ability-shoppinglist/pkg/shoppinglist"
)

var shoppingListClient *shoppinglist.Client

const (
	deleteAction="DELETE"
	addAction="ADD"
	itemsSlot="ITEMS"
	itemEntity="SHOPITEM"
)

// fun main()
func main() {
	// Read configuration
	conf := ability.ReadConfiguration()

	// Initialize client for shopping list tool
	shoppingListClient = shoppinglist.NewClient(conf.Tools["shoppinglist"].Host, conf.Tools["shoppinglist"].Port)

	// Initialize server
	server := ability.NewServer("Shopping List Ability", conf.Server.Port)
	server.RegisterIntentRule("TRIGGER_SHOPPING_LIST", triggerShoppingListHandler)
	server.RegisterIntentRule("REMOVE_FROM_LIST", removeFromListHandler)
	server.RegisterIntentRule("REMOVE_FROM_SHOPPING_LIST", removeFromListHandler)
	server.RegisterIntentRule("ADD_TO_LIST", addToListHandler)
	server.RegisterIntentRule("ADD_TO_SHOPPING_LIST", addToListHandler)
	server.RegisterIntentRule("EMPTY_LIST_ITEMS", emptyShoppingListHandler)
	server.RegisterIntentRule("COUNT_LIST_ITEMS", countShoppingListHandler)
	server.RegisterIntentRule("LIST_LIST_ITEMS", listShoppingListHandler)
	server.RegisterRule(isRemoveContext, removeFromListHandler)
	server.RegisterRule(isAddContext, addToListHandler)
	server.Serve()
}

func isRemoveContext(req *ability.Request) bool {
	return req.IsInSlotFillingAction(deleteAction)
}

func isAddContext(req *ability.Request) bool {
	return req.IsInSlotFillingAction(addAction)
}

func removeFromListHandler(req *ability.Request, resp *ability.Response) {
	// Retrieve only shopping items from NLU entities
	items := collectItemsFromRequest(req)

	// If we don't find any items in the request, we ask to user
	if len(items) == 0 {
		resp.Nlg.Sentence = "What do you want to delete from your shopping list ?"
		resp.Context.SlotFilling.Action = deleteAction
		resp.Context.SlotFilling.MissingSlots = []string{itemsSlot}
		resp.AutoReprompt = true
		return
	}

	// Remove these items from the shopping list
	if err := shoppingListClient.RemoveItems(items); err != nil {
		resp.Nlg.Sentence = "Error removing item from your shopping list."
		return
	}

	// Build the NLG answer
	resp.Nlg.Sentence = "I removed {{items}} from your shopping list."
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func addToListHandler(req *ability.Request, resp *ability.Response) {
	// Retrieve only shopping items from NLU entities
	items := collectItemsFromRequest(req)

	// If we don't find any items in the request, we ask to user
	if len(items) == 0 {
		resp.Nlg.Sentence = "What do you want to add to your shopping list ?"
		resp.Context.SlotFilling.Action = addAction
		resp.Context.SlotFilling.MissingSlots = []string{itemsSlot}
		resp.AutoReprompt = true
		return
	}

	// Add these items into the shopping list
	if err := shoppingListClient.AddItems(items); err != nil {
		resp.Nlg.Sentence = "Error adding item to your shopping list."
		return
	}

	// Build the NLG answer
	resp.Nlg.Sentence = "I added {{items}} to your shopping list"
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func triggerShoppingListHandler(_ *ability.Request, resp *ability.Response) {
	items, err := shoppingListClient.GetItems()
	if err != nil {
		resp.Nlg.Sentence = "Error receiving items from your shopping list."
		return
	}
	// Build the NLG answer
	resp.Nlg.Sentence = "You have {{count}} items in your main shopping list, what do you want to do ?"
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "count",
		Value: len(items),
		Type:  "string",
	}}
	resp.AutoReprompt = true
}

func emptyShoppingListHandler(_ *ability.Request, resp *ability.Response) {
	// Here we count the items
	var count = -1
	items, err := shoppingListClient.GetItems()
	if err != nil {
		count = len(items)
	}

	// If there is no item, we don't go further
	if count == 0 {
		resp.Nlg.Sentence = "Your shopping list is already empty."
		return
	}

	err = shoppingListClient.RemoveAllItems()
	if err != nil {
		resp.Nlg.Sentence = "Error removing all items from your shopping list."
		return
	}

	// Build the NLG answer
	if count < 0 {
		resp.Nlg.Sentence = "Your shopping list has been cleared."
		return
	}

	resp.Nlg.Sentence = "{{count}} items has been removed from your shopping list."
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}}
}

func countShoppingListHandler(_ *ability.Request, resp *ability.Response) {
	items, err := shoppingListClient.GetItems()
	if err != nil {
		resp.Nlg.Sentence = "Error counting items from your shopping list."
		return
	}
	count := len(items)
	if count <= 0 {
		resp.Nlg.Sentence = "You don't have any elements in your shopping list."
		return
	}

	resp.Nlg.Sentence = "You have {{count}} items in your shopping list."
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}}
}

func listShoppingListHandler(_ *ability.Request, resp *ability.Response) {
	items, err := shoppingListClient.GetItems()
	if err != nil {
		resp.Nlg.Sentence = "Error receiving items from your shopping list."
		return
	}
	count := len(items)
	if count <= 0 {
		resp.Nlg.Sentence = "You don't have any elements in your shopping list."
		return
	}

	if count == 1 {
		resp.Nlg.Sentence = "You only have one element in your shopping list. There is {{item}."
		resp.Nlg.Params = []ability.NLGParam{{
			Name:  "item",
			Value: items[0],
			Type:  "string",
		}}
		return
	}

	resp.Nlg.Sentence = "You have {{count}} items in your shopping list. There are {{items}}."
	resp.Nlg.Params = []ability.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}, {
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func collectItemsFromRequest(req *ability.Request) []string {
	return req.GetEntitiesByLabel(itemEntity)
}
