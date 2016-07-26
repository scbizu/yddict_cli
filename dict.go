package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mkideal/cli"
)

type Dict struct {
	// Error-code
	ErrorCode int
	// Query string. Could be different from the request.
	Query string
	// Translations
	Translation []string
	// Basic dictionary result. Could be <nil>
	Basic *struct {
		// Phonetic data. Could be ""
		Phonetic string
		// Explains
		Explains []string
	}
	// Web mining dictionary result. Could be a zero-length slice
	Web []struct {
		// Entry key
		Key string
		// Explains
		Value []string
	}
}

type argT struct {
	cli.Helper
	Word string `cli:"*word" usage:"input the English/Chinese word...This is required"`
	Info string `cli:"info" usage:"input the content information.Note:less for quick search and show less info,more for more translations" dft:"less"`
}

func main() {
	dict := &Dict{}
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		//urlencode the value
		urlword := url.QueryEscape(argv.Word)
		url := "http://fanyi.youdao.com/openapi.do?keyfrom=LocalDict&key=798282145&type=data&doctype=json&version=1.1&q=" + urlword
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err.Error())
		}
		Derr := json.NewDecoder(resp.Body).Decode(dict)
		if Derr != nil {
			log.Fatal(err.Error())
		}
		if argv.Info == "more" {
			ctx.JSONln(dict.Basic)
		} else {
			fmt.Println("result:", dict.Translation)
		}

		defer resp.Body.Close()
		return nil
	}, "cli for my dict based on youdao")
}
