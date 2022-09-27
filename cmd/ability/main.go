package main

import (
	"github.com/milobella/ability-sdk-go/pkg/config"
	"github.com/milobella/ability-sdk-go/pkg/model"
	"github.com/milobella/ability-sdk-go/pkg/server"
	"github.com/milobella/ability-sdk-go/pkg/server/conditions"
	"github.com/milobella/ability-shoppinglist/pkg/shoppinglist"
)

var shoppingListClient *shoppinglist.Client

const (
	deleteAction = "DELETE"
	addAction    = "ADD"
	itemsSlot    = "ITEMS"
	itemEntity   = "SHOPITEM"
)

// fun main()
func main() {
	// Read configuration
	conf := config.Read()

	// Initialize client for shopping list tool
	shoppingListClient = shoppinglist.NewClient(conf.Tools["shoppinglist"].Host, conf.Tools["shoppinglist"].Port)

	// Initialize server
	srv := server.New("Shopping List", conf.Server.Port)

	// Register first the conditions on actions because they have priority on intents.
	// The condition returns true if an action is pending.
	srv.Register(conditions.IfInSlotFilling(deleteAction), removeFromListHandler)
	srv.Register(conditions.IfInSlotFilling(addAction), addToListHandler)

	// Then we register intents routing rules.
	// It means that if no pending action has been found in the context, we'll use intent to decide the handler.
	srv.Register(conditions.IfIntents("TRIGGER_SHOPPING_LIST"), triggerShoppingListHandler)
	srv.Register(conditions.IfIntents("REMOVE_FROM_LIST", "REMOVE_FROM_SHOPPING_LIST"), removeFromListHandler)
	srv.Register(conditions.IfIntents("ADD_TO_LIST", "ADD_TO_SHOPPING_LIST"), addToListHandler)
	srv.Register(conditions.IfIntents("EMPTY_LIST_ITEMS"), emptyShoppingListHandler)
	srv.Register(conditions.IfIntents("COUNT_LIST_ITEMS"), countShoppingListHandler)
	srv.Register(conditions.IfIntents("LIST_LIST_ITEMS"), listShoppingListHandler)

	srv.Serve()
}

func removeFromListHandler(req *model.Request, resp *model.Response) {
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
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func addToListHandler(req *model.Request, resp *model.Response) {
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
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func triggerShoppingListHandler(_ *model.Request, resp *model.Response) {
	items, err := shoppingListClient.GetItems()
	if err != nil {
		resp.Nlg.Sentence = "Error receiving items from your shopping list."
		return
	}
	// Build the NLG answer
	resp.Nlg.Sentence = "You have {{count}} items in your main shopping list, what do you want to do ?"
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "count",
		Value: len(items),
		Type:  "string",
	}}
	resp.AutoReprompt = true
}

func emptyShoppingListHandler(_ *model.Request, resp *model.Response) {
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
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}}
}

func countShoppingListHandler(_ *model.Request, resp *model.Response) {
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
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}}
}

func listShoppingListHandler(_ *model.Request, resp *model.Response) {
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
		resp.Nlg.Params = []model.NLGParam{{
			Name:  "item",
			Value: items[0],
			Type:  "string",
		}}
		return
	}

	resp.Nlg.Sentence = "You have {{count}} items in your shopping list. There are {{items}}."
	resp.Nlg.Params = []model.NLGParam{{
		Name:  "count",
		Value: count,
		Type:  "string",
	}, {
		Name:  "items",
		Value: items,
		Type:  "enumerated_list",
	}}
}

func collectItemsFromRequest(req *model.Request) []string {
	return req.GetEntitiesByLabel(itemEntity)
}
