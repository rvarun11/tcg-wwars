package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/shomali11/slacker"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"tcg.ai/tool"
)

type Query struct {
	Body string `json:"body"`
}

type Reply struct {
	Body string `json:"body"`
}

func printCmdEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Event")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println("-----------")
	}
}

func main() {
	serviceUrl := "http://localhost:8000/"
	tool.SetTokens()
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	//go printCmdEvents(bot.CommandEvents())

	defaultCmd := func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Hi!! I'm sorry I didn't get you. " +
			"Please type *ask* followed by your query and I'll see what I can do. :) \n" +
			"For eg: _ask How far is the sun?_")
	}

	bot.DefaultCommand(defaultCmd)
	//TODO: 1. find a way to get query witout ask. 2. handle edge cases with invalid input.
	bot.Command("ask <query>", &slacker.CommandDefinition{
		Description: "ask Ask a question",
		Example:     "ask How far is the sun?",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("query")
			reply := getReply(serviceUrl, query)
			response.Reply(reply)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func getReply(serviceUrl string, q string) string {
	fmt.Println("Connecting to FastAPI at: ", serviceUrl)
	query := &Query{Body: q}
	jsonBytes, _ := json.Marshal(query)

	res, err := http.Post(serviceUrl, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println("Response Status: ", res.Status)
	//fmt.Println("Response Headers: ", res.Header)
	//fmt.Println("Response Body: ", string(body))

	reply := Reply{}
	_ = json.Unmarshal(body, &reply)

	return reply.Body
}
