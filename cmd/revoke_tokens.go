//go:build !production

package cmd

import (
	"context"
	"fmt"

	"github.com/Binaretech/classroom-auth/cache"
	"github.com/spf13/cobra"
)

var revokeTokens = &cobra.Command{
	Use:   "revoke:all",
	Short: "revoke all tokens",
	Run: func(cmd *cobra.Command, args []string) {
		iterator := cache.Scan(context.Background(), 0, "*", 0).Iterator()
		for iterator.Next(context.Background()) {
			fmt.Println(iterator.Val())
		}
	},
}

func init() {
	rootCmd.AddCommand(revokeTokens)
}
