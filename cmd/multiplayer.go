/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

// multiplayerHostCmd represents the multiplayer host command
var multiplayerHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Host a typistone race",
	Long: `Allows you to create a typistone race and returns an invite link.
Other players can then join using this link.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("multiplayer host called")
	},
}

func init() {
	multiplayerCmd.AddCommand(multiplayerHostCmd)
	rootCmd.AddCommand(multiplayerCmd)

	multiplayerHostCmd.Flags().Int("game-size", 2, "Specify the number of players that will join this game")
	multiplayerCmd.Flags().Bool("guest", false, "Play as guest")
	multiplayerCmd.Flags().String("link", "localhost:8080", "Join a game using an invite link")
}
