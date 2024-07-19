package cmd

import (
	"github.com/gatariee/sanity/internal/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "prepare service files",
	Long:  `prep service files without having to worry about accidentally leaving a flag in there`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		ff, _ := cmd.Flags().GetString("ff")
		exclude, _ := cmd.Flags().GetString("exclude")

		err := service.PrepareService(input, ff, exclude)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.Flags().String("input", "", "path to service file")
	serviceCmd.MarkFlagRequired("input")

	serviceCmd.Flags().String("ff", "", "flag format, only prefix") // i.e if FLAG{...} then FLAG
	serviceCmd.MarkFlagRequired("ff")

	serviceCmd.Flags().Bool("zip", false, "zip dist folder")

	/* useful for known edge cases */
	serviceCmd.Flags().String("exclude", "", "exclude from flag check") // i.e FLAG{THIS_IS_A_FAKE_FLAG}
}