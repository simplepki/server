package main

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var rootCmd = &cobra.Command{
	Use:   "simplepkid",
	Short: "simplepkid is server cli",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs the server using the settings.json file located in: [/etc/simplepki/, /opt/simplepki/ $HOME/.simplepki/, ./]",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting simplepkid")
	},
}
