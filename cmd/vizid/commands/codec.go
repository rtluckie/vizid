package commands

import (
	"fmt"

	"github.com/ryanl/vizid/internal/codec"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode <vizid>",
	Short: "Decode a VIZID into its ASCII wire format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ascii, err := codec.DecodeVIZToASCII(args[0])
		if err != nil {
			return err
		}
		fmt.Println(ascii)
		return nil
	},
}

var encodeCmd = &cobra.Command{
	Use:   "encode <ascii>",
	Short: "Encode an ASCII wire ID into VIZ glyph form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viz, err := codec.EncodeASCIIToVIZ(args[0])
		if err != nil {
			return err
		}
		fmt.Println(viz)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(encodeCmd)
}
