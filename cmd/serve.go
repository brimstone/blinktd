package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brimstone/blinkt_go/sysfs"
	"github.com/brimstone/blinktd/types"
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

var pixels [8]types.Pixel

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the blinktd server",
	Long: `Run a server that manages the pixels on your Blinktd!
Use the client or API to change the value of the pixels.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := httpd.New(httpd.Port(8000))
		if err != nil {
			return err
		}
		h.HandleFunc("/v1/led", handleLed)

		blinkt.SetClearOnExit(true)

		blinkt.Setup()

		time.Sleep(100 * time.Millisecond)
		blinkt.Clear()
		pixels[0] = types.Pixel{
			Green:  255,
			Format: types.PixelMorse,
			Value:  0,
		}
		for id := range pixels {
			go func(id int) {
				for {
					pixel := pixels[id]
					if pixel.Red == 0 && pixel.Green == 0 && pixel.Blue == 0 {
						continue
					}
					if pixel.Format == types.PixelMorse {
						morsePixel(pixel.Value, id, pixel.Red, pixel.Green, pixel.Blue)
					} else {
						blinkt.SetPixel(id, pixel.Red, pixel.Green, pixel.Blue)
						time.Sleep(time.Second)
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
		if dot == 0 {
			time.Sleep(dotlen * 3)
		} else {
			time.Sleep(dotlen)
		}
		blinkt.SetPixel(pixel, 0, 0, 0)
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
	var requestPixel types.Pixel
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
