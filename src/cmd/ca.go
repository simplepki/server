package main

import (
	"log"
	"strings"

	"github.com/jtaylorcpp/piv-go/piv"
	"github.com/simplepki/server/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(caCmd)
	caCmd.AddCommand(yubiCmd, caInitCmd)
	yubiCmd.AddCommand(listYubiCmd)
}

var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "ca tools",
}

var caInitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize CA",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("initializing ca")
		if !config.IsCAEnabled() {
			log.Println("no ca definition found")
			return
		}
	},
}

var yubiCmd = &cobra.Command{
	Use:   "yubikey",
	Short: "yubikey related tools",
}

var listYubiCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all connected yubikeys",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("getting all attacked yubikeys")
		cards, err := piv.Cards()
		if err != nil {
			log.Printf("no yubikey present w/ error: %s\n", err.Error())
		}

		for _, card := range cards {
			if strings.Contains(strings.ToLower(card), "yubikey") {
				yk, err := piv.Open(card)
				if err != nil {
					log.Printf("unable to open yubikey: %s\n", cards)
					continue
				}

				serial, err := yk.Serial()
				if err != nil {
					log.Printf("unable to get yubikey serial number: %v\n", serial)
					continue
				}
				log.Printf("SN: %v, Name: %v\n", serial, card)
			}
		}
	},
}
