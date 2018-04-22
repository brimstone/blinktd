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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
