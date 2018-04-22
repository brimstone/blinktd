package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/alexellis/blinkt_go/sysfs"
	httpd "github.com/brimstone/go-httpd"
	"github.com/spf13/cobra"
)

var blinkt = sysfs.NewBlinkt(0.2)
var dotlen = time.Second / 10

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

type PixelFormat string

var PixelMorse PixelFormat = "morse"
var PixelSolid PixelFormat = "solid"

type Pixel struct {
	ID     int         `json:id`
	Red    int         `json:red`
	Green  int         `json:green`
	Blue   int         `json:blue`
	Format PixelFormat `json:format`
	Value  int64       `json:value`
}

var pixels [8]Pixel

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
		pixels[0] = Pixel{
			Green:  255,
			Format: PixelMorse,
			Value:  0,
		}
		for id := range pixels {
			go func(id int) {
				for {
					pixel := pixels[id]
					if pixel.Format == PixelMorse {
						morsePixel(pixel.Value, id, pixel.Red, pixel.Green, pixel.Blue)
					} else {
						blinkt.SetPixel(id, pixel.Red, pixel.Green, pixel.Blue)
						blinkt.Show()
						sysfs.Delay(100)
					}
				}
			}(id)
		}

		return h.ListenAndServe()

	},
}

func morsePixel(value int64, pixel int, r int, g int, b int) {
	for _, dot := range morseDigit[value] {
		blinkt.SetPixel(pixel, r, g, b)
		blinkt.Show()
		if dot == 0 {
			time.Sleep(dotlen * 3)
		} else {
			time.Sleep(dotlen)
		}
		blinkt.SetPixel(pixel, 0, 0, 0)
		blinkt.Show()
		time.Sleep(dotlen)
	}
	time.Sleep(dotlen * 3)
}

func handleLed(w http.ResponseWriter, r *http.Request) {
	// TODO handle auth

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()
	var requestPixel Pixel
	err = json.Unmarshal(body, &requestPixel)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	pixels[requestPixel.ID] = requestPixel
	w.Write([]byte("OK"))
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
