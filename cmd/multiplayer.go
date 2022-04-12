/*
Copyright © 2022
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"

	serverModels "github.com/DelusionalOptimist/typistone-server/models"
	"github.com/DelusionalOptimist/typistone/models"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	host string
	lobbySize int
	lobbyID string
)

// multiplayerCmd represents the multiplayer command
var multiplayerCmd = &cobra.Command{
	Use:   "multiplayer",
	Short: "Play typistone with other players",
	Long: `Play typistone with other players across the world.
You can host a game yourself, join someone's game with a link or get automatically matched against players.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("multiplayer called")
	},
}

// multiplayerCreateCmd represents the multiplayer create command
var multiplayerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new lobby",
	Long: `Allows you to create a lobby and returns an invite link.
Other players can then join the lobby using this link.`,
	Run: func(cmd *cobra.Command, args []string) {
		// opens up a websocket connection with the host
		u := url.URL{ Scheme: "ws", Host: host, Path: "/create" }
		fmt.Println("Connecting to ", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial: ", err)
		}
		defer c.Close()

		// create a new lobby and send its config to server
		lobby := &serverModels.Lobby{}
		lobby.LobbyConfig.LobbySize = lobbySize
		err = c.WriteJSON(lobby)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Read new lobby details received
		err = c.ReadJSON(lobby)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Lobby ID: ", lobby.LobbyID)

		// Read the player data
		player := &serverModels.Player{}
		err = c.ReadJSON(player)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Player ID: ", player.PlayerID)

		// wait till receiving the "Starting game" msg
		var msg string
		for msg != "Starting game" {
			_, m, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			msg = string(m)
		}
		fmt.Println(msg)

		// text to type received from server
		_, text, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// create a new game model and receive in it data sent from the
		// central server
		gameData := &models.Game{}
		gameData.Text = string(text)

		gameDataReceived := make(chan models.Game)
		go func() {
			for {
				err = c.ReadJSON(gameData)
				if err != nil {
					fmt.Println(err)
				}
				gameDataReceived <- *gameData
			}
		}()

		//gameDataSent := make(chan models.Game)
		//go func(){
		//	gameData := <- gameDataSent
		//	err = c.WriteJSON(gameData)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}()

		// do stuff
		for {
			select {
			case data := <- gameDataReceived:
				fmt.Println(data)
			}
		}
	},
}

// multiplayerJoinCmd represents the multiplayer create command
var multiplayerJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "join a lobby with the given ID",
	Long: `Allows you to join a lobby.`,
	Run: func(cmd *cobra.Command, args []string) {
		// opens up a websocket connection with the host
		u := url.URL{Scheme: "ws", Host: host, Path: fmt.Sprint("/join/", lobbyID)}
		log.Printf("Connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial: ", err)
		}
		defer c.Close()

		// Read the player data
		player := &serverModels.Player{}
		err = c.ReadJSON(player)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Player ID: ", player.PlayerID)

		// wait till receiving the "Starting game" msg
		var msg string
		for msg != "Starting game" {
			_, m, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			msg = string(m)
		}
		fmt.Println(msg)

		// text to type received from server
		_, text, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		gameData := &models.Game{}
		gameData.Text = string(text)

		var fakePercentage = 10.0
		for {
			gameData.Percentages[player.PlayerID-1] = fakePercentage
			fakePercentage++
			err := c.WriteJSON(gameData)
			if err != nil {
				log.Println(err)
				return
			}
		}

	},
}


func init() {
	multiplayerCreateCmd.Flags().IntVar(&lobbySize, "lobby-size", 2, "Specify the number of players that will join this game")
	multiplayerCmd.Flags().Bool("guest", false, "Play as guest")
	multiplayerCmd.Flags().StringVar(&host, "host", "localhost:8080", "URL of the server which hosts the lobby")
	multiplayerJoinCmd.Flags().StringVar(&lobbyID, "lobby-id", "", "ID of the lobby to join")

	multiplayerCmd.AddCommand(multiplayerCreateCmd)
	multiplayerCmd.AddCommand(multiplayerJoinCmd)
	rootCmd.AddCommand(multiplayerCmd)
}
