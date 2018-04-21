package cmd

import (
	"fmt"

	"github.com/alexellis/blinkt_go/sysfs"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		brightness := 1.0
		blinkt := sysfs.NewBlinkt(brightness)

		blinkt.SetClearOnExit(true)

		blinkt.Setup()

		sysfs.Delay(100)

		r := 0
		g := 0
		b := 255

		for {
			for pixel := 0; pixel < 8; pixel++ {
				blinkt.Clear()
				blinkt.SetPixel(pixel, r, g, b)
				blinkt.Show()
				sysfs.Delay(100)
			}
			for pixel := 7; pixel > 0; pixel-- {
				blinkt.Clear()
				blinkt.SetPixel(pixel, r, g, b)
				blinkt.Show()
				sysfs.Delay(100)
			}
		}
		blinkt.Clear()
		blinkt.Show()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
