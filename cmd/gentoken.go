package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/brimstone/blinktd/types"
	"github.com/brimstone/jwt/jwt"
	"github.com/spf13/cobra"
)

// genTokenCmd represents the gentoken command
var genTokenCmd = &cobra.Command{
	Use:   "gentoken",
	Short: "Generate a token to control authorization",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			return errors.New("Must specify a key with -k")
		}

		authToken := types.AuthToken{}

		authToken.Pixels, err = cmd.Flags().GetIntSlice("pixel")

		if len(authToken.Pixels) == 0 {
			return errors.New("Must specify some pixels")
		}

		authTokenJSON, err := json.Marshal(authToken)
		if err != nil {
			return err
		}

		log.Println(string(authTokenJSON))

		token, err := jwt.GenToken(key, authTokenJSON)
		if err != nil {
			return err
		}
		fmt.Println(token)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genTokenCmd)
	genTokenCmd.Flags().StringP("key", "k", "", "The secret key used to sign the token.")
	genTokenCmd.Flags().IntSliceP("pixel", "p", []int{}, "Allowed Pixel")
}
