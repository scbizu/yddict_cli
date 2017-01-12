package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
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

		//if the word is checked,then put it into the db
		db, dberr := bolt.Open("mydict.db", 0600, nil)
		defer db.Close()
		if dberr != nil {
			color.Red("db create failed")
		}
		//create bucket
		bterr := db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte("yddict"))
			if err != nil {
				return fmt.Errorf("Cannot Create Bucket:%s", err)
			} else {
				return nil
			}
		})
		if bterr != nil {
			color.Red("Create/Access bucket failed")
		}

		//put word into dict
		db.Update(func(tx *bolt.Tx) error {
			bt := tx.Bucket([]byte("yddict"))
			value := bt.Get([]byte(urlword))
			if value == nil {
				err = bt.Put([]byte(urlword), []byte("checked"))
				color.Green("remembered!")
				return err
			} else {
				color.Blue("You have checked it")
				//	err := bt.Put([]byte(urlword), []byte("1"))
				return nil
			}

		})

		if argv.Info == "more" {
			ctx.JSONln(dict.Basic)
		} else {
			fmt.Println("result:", dict.Translation)
		}

		defer resp.Body.Close()
		return nil
	}, "cli for my dict based on youdao")
}
