package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/celian-garcia/gonfig"
	"github.com/sirupsen/logrus"
	"gitlab.milobella.com/milobella/ability-sdk-go/pkg/ability"
	"gitlab.milobella.com/milobella/shoppinglist-ability/internal/config"
	"gitlab.milobella.com/milobella/shoppinglist-ability/pkg/shoppinglist"
)

var shoppingListClient *shoppinglist.Client

type Configuration struct {
	Server     config.ServerConfiguration
	Tool       config.ToolConfiguration
	ConfigFile string `short:"c"`
}

func (c Configuration) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Configuration serialization error %v", err)
	}
	return string(b)
}

var conf *Configuration

var deleteAction string
var addAction string
var itemsSlot string
var itemEntity string


//TODO: use this init function to initialize variables instead of initialize on top
func init() {
	deleteAction = "DELETE"
	addAction = "ADD"
	itemsSlot = "ITEMS"
	itemEntity = "SHOPITEM"

	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)
}

// fun main()
func main() {
	//TODO: change the configuration for vyper
	conf = &Configuration{}

	// Load the configuration from file or parameter or env
	err := gonfig.Load(conf, gonfig.Conf{
		ConfigFileVariable: "configfile", // enables passing --configfile myfile.conf

		FileDefaultFilename: "config/ability.toml",
		FileDecoder:         gonfig.DecoderTOML,

		EnvPrefix: "ABILITY_",
	})

	if err != nil {
		logrus.Fatalf("Error reading config : %s", err)
	} else {
		logrus.Infof("Successfully readen configuration file : %s", conf.ConfigFile)
		logrus.Debugf("-> %+v", conf)
	}

	// Initialize client for shopping list tool
	shoppingListClient = shoppinglist.NewClient(conf.Tool.Host, conf.Tool.Port)

	// Initialize server
	server := ability.NewServer("Shopping List Ability", conf.Server.Port)
	// TODO: remove it to use only rule
	if err = server.RegisterIntent("ADD_TO_LIST", func(req ability.Request, resp *ability.Response) {
		addToListHandler(&req, resp)
	}); err != nil {
		logrus.Errorf(err.Error())
	}
	// TODO: remove it to use only rule
	if err = server.RegisterIntent("TRIGGER_SHOPPING_LIST", func(req ability.Request, resp *ability.Response) {
		triggerShoppingListHandler(&req, resp)
	}); err != nil {
		logrus.Errorf(err.Error())
	}
	server.RegisterIntentRule("REMOVE_FROM_LIST", removeFromListHandler)
	server.RegisterIntentRule("ADD_TO_LIST", addToListHandler)
	server.RegisterIntentRule("TRIGGER_SHOPPING_LIST", triggerShoppingListHandler)
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

func triggerShoppingListHandler(req *ability.Request, resp *ability.Response) {
	items, err := shoppingListClient.GetItems()
	if err != nil {
		resp.Nlg.Sentence = "Error receiving item from your shopping list."
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

func collectItemsFromRequest(req *ability.Request) []string {
	return req.GetEntitiesByLabel(itemEntity)
}
