package helpers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var aid string = "U01MENEF744"
var rid string = "U01MM6PS3PB"

/*
	after getting all files:
	1. read the json files
		1.1 if reply_users contains alex_uid & rahul_id
			1.1.2 parse replies [] and retrieve the timestamp.
			1.1.3 print the ts for now. ->
*/

type Message []struct {
	Type           string   `json:"type"`
	Text           string   `json:"text"`
	User           string   `json:"user"`
	ReplyUserCount int      `json:"reply_user_count"`
	ReplyUsers     []string `json:"reply_users"`
	Replies        Replies  `json:"replies"`
}

type Replies []struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

func walk(path string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		fmt.Println(path)
	}
	return nil
}

// RetrieveTS will retrieve all TS for alex and rahul's reply.
func RetrieveTS(path string) {
	bytes, _ := os.ReadFile(path)
	var res Message
	json.Unmarshal(bytes, &res)
	fmt.Println(res[0].Replies[0].Ts)
}

func getFiles() {
	filepath.WalkDir("./data", walk)
}
