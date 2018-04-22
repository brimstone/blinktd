package cmd

import (
	"net/http"

	"github.com/alexellis/blinkt_go/sysfs"
	httpd "github.com/brimstone/go-httpd"
	"github.com/spf13/cobra"
)

var blinkt = sysfs.NewBlinkt(1.0)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := httpd.New(httpd.Port(8000))
		if err != nil {
			return err
		}
		h.HandleFunc("/v1/led", handleLed)

		blinkt.SetClearOnExit(true)

		blinkt.Setup()

		sysfs.Delay(100)
		blinkt.Clear()
		blinkt.SetPixel(0, 0, 255, 0)
		blinkt.Show()

		return h.ListenAndServe()

	},
}

func handleLed(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
	blinkt.Clear()
	blinkt.SetPixel(0, 255, 0, 0)
	blinkt.Show()
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
