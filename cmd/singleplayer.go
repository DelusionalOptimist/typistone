/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/typer/pkg/model"
	wrap "github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

const (
	blue          = "#4776E6"
	purple        = "#8E54E9"
	words         = 15
	defaultWidth  = 60
	DefaultLength = 20
	MaxLength     = 500
	author = "william+shakespeare"
)

type response struct {
	CurrentPage int `json:"current_page"`
	TotalPages int `json:"total_pages"`
	Quotes []quote `json:"quotes"`
}

type quote struct {
	Quote string `json:"quote"`
	Author string `json:"author"`
	Publication string `json:"publication"`
}

// singleplayerCmd represents the singleplayer command
var singleplayerCmd = &cobra.Command{
	Use:   "singleplayer",
	Short: "Practice typing with typistone by yourself",
	Run: func(cmd *cobra.Command, args []string) {
		// do bubble tea stuff

		// TODO: improve fetching logic
		rand.Seed(time.Now().UnixNano())
		pageIndex := rand.Intn(100)
		quotesAPIURL := fmt.Sprint("https://goodquotesapi.herokuapp.com/author/" + author + "/?page=" + fmt.Sprint(pageIndex))

		req, err := http.Get(quotesAPIURL)
		//req, err := http.
		if err != nil {
			log.Fatalln(err)
		}

		body := json.NewDecoder(req.Body)
		obj := &response{}
		err = body.Decode(obj)
		if err != nil {
			log.Fatalln(err)
		}

		quoteNo := rand.Intn(20)
		text := obj.Quotes[quoteNo].Quote

		// stuff required
		bar, err := progress.NewModel(progress.WithScaledGradient(blue, purple))
		if err != nil {
			log.Fatalln(err)
		}

		program := tea.NewProgram(model.Model{
			Progress: bar,
			Text:     wrap.WrapString(text, defaultWidth),
		})

		program.Start()
	},
}

func init() {
	rootCmd.AddCommand(singleplayerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// singleplayerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// singleplayerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
