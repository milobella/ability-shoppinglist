package main

import (
    "encoding/json"
    "github.com/juju/loggo"
    "github.com/stevenroose/gonfig"
    "gitlab.milobella.com/milobella/ability-sdk-go/pkg/ability"
    "gitlab.milobella.com/milobella/shoppinglist-ability/internal/config"
    "gitlab.milobella.com/milobella/shoppinglist-ability/pkg/shoppinglist"
    "gitlab.milobella.com/milobella/oratio/pkg/anima"
    "log"
)

var shoppingListClient = shoppinglist.NewClient("http://0.0.0.0", 4848)

type Configuration struct {
    Server     config.ServerConfiguration
    ConfigFile string `short:"c"`
}

func (c Configuration) String() string {
    b, err := json.Marshal(c)
    if err != nil {
        log.Fatal("Configuration serialization error %s", err)
    }
    return string(b)
}

var conf *Configuration

// fun main()
func main() {
    conf = &Configuration{}

    // Load the configuration from file or parameter or env
    err := gonfig.Load(conf, gonfig.Conf{
        ConfigFileVariable: "configfile", // enables passing --configfile myfile.conf

        FileDefaultFilename: "config/shoppinglist-ability.toml",
        FileDecoder:         gonfig.DecoderTOML,

        EnvPrefix: "ABILITY_",
    })

    logger := loggo.GetLogger("shoppinglist-ability.main")
    if err != nil {
        loggo.ConfigureLoggers("<root>=INFO")
        logger.Criticalf("Error reading config : %s", err)
    } else {
        loggo.ConfigureLoggers(conf.Server.LogLevel)
        logger.Infof("Successfully readen configuration file : %s", conf.ConfigFile)
        logger.Debugf("-> %+v", conf)
    }

    server := ability.NewServer(conf.Server.Port)
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
