package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/brimstone/blinktd/types"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a pixel",
	Long:  `Connect to the daemon run by serve and set a pixel.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// First ID
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		if id < 0 || id > 7 {
			return errors.New("ID value is out of range")
		}

		// Create pixel object.
		pixel := types.Pixel{
			ID: id,
		}

		pixel.Red, err = cmd.Flags().GetInt("red")
		if err != nil {
			return err
		}
		if pixel.Red < 0 || pixel.Red > 255 {
			return errors.New("Red is out of range")
		}

		// Green
		pixel.Green, err = cmd.Flags().GetInt("green")
		if err != nil {
			return err
		}
		if pixel.Green < 0 || pixel.Green > 255 {
			return errors.New("Green is out of range")
		}

		// Blue
		pixel.Blue, err = cmd.Flags().GetInt("blue")
		if err != nil {
			return err
		}
		if pixel.Blue < 0 || pixel.Blue > 255 {
			return errors.New("Blue is out of range")
		}

		// Morse?
		morse, err := cmd.Flags().GetBool("morse")
		if err != nil {
			return err
		}
		if morse {
			pixel.Format = types.PixelMorse
		}

		// Value
		value, err := cmd.Flags().GetInt("value")
		pixel.Value = int64(value)
		if err != nil {
			return err
		}
		if pixel.Value < 0 {
			return errors.New("Pixel value is too low")
		}
		if pixel.Value > 9 {
			return errors.New("Pixel value is too high")
		}

		// Server
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}

		// Port
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return err
		}

		pixelJSON, err := json.Marshal(pixel)
		body := bytes.NewReader(pixelJSON)
		req, err := http.NewRequest("POST", "http://"+server+":"+port+"/v1/led", body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return errors.New(string(body))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("server", "s", "localhost", "Server address")
	setCmd.Flags().StringP("port", "p", "8000", "Server Port")

	setCmd.Flags().IntP("id", "i", 0, "Pixel id, 0-7")
	setCmd.Flags().IntP("red", "r", 0, "Amount of red, 0-255")
	setCmd.Flags().IntP("green", "g", 0, "Amount of green, 0-255")
	setCmd.Flags().IntP("blue", "b", 0, "Amount of blue, 0-255")
	setCmd.Flags().BoolP("morse", "m", false, "If morse code should be used to show the value")
	setCmd.Flags().IntP("value", "v", 0, "Value to show with morse code")

	setCmd.Flags().StringP("token", "t", "", "Token for authorization")
}
