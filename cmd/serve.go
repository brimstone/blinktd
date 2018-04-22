package cmd

import (
	"net/http"
	"time"

	"github.com/alexellis/blinkt_go/sysfs"
	httpd "github.com/brimstone/go-httpd"
	"github.com/spf13/cobra"
)

var blinkt = sysfs.NewBlinkt(1.0)
var dotlen = time.Second / 2
var status = int64(0)

var morseDigit = [][]int{
	[]int{0, 0, 0, 0, 0},
	[]int{1, 0, 0, 0, 0},
	[]int{1, 1, 0, 0, 0},
	[]int{1, 1, 1, 0, 0},
	[]int{1, 1, 1, 1, 0},
	[]int{1, 1, 1, 1, 1},
	[]int{0, 1, 1, 1, 1},
	[]int{0, 0, 1, 1, 1},
	[]int{0, 0, 0, 1, 1},
	[]int{0, 0, 0, 0, 1},
}

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
		go func() {
			for {
				for _, dot := range morseDigit[status] {
					blinkt.SetPixel(0, 0, 255, 0)
					blinkt.Show()
					time.Sleep(dotlen)
					if dot == 0 {
						time.Sleep(dotlen)
						time.Sleep(dotlen)
					}
					blinkt.SetPixel(0, 0, 0, 0)
					blinkt.Show()
					time.Sleep(dotlen)
				}
			}
		}()

		return h.ListenAndServe()

	},
}

func handleLed(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
	blinkt.Clear()
	blinkt.SetPixel(1, 255, 0, 0)
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
