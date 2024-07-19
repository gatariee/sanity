package cmd

import (
	"github.com/gatariee/sanity/internal/service"
	"github.com/gatariee/sanity/internal/utility"
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
		batch, _ := cmd.Flags().GetBool("batch")
		temp_path, _ := cmd.Flags().GetString("name")

		err := service.PrepareService(input, ff, exclude, batch, temp_path)
		if err != nil {
			panic(err)
		}

		zip, _ := cmd.Flags().GetString("zip")
		if zip != "" {
			err = utility.ZipFiles(temp_path, zip)
			if err != nil {
				panic(err)
			}

			cleanup, _ := cmd.Flags().GetBool("cleanup")
			if cleanup {
				err = utility.RemoveFile(temp_path)
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.Flags().String("name", "", "name of the folder to be created")
	serviceCmd.MarkFlagRequired("name")

	serviceCmd.Flags().String("input", "", "path to service file")
	serviceCmd.MarkFlagRequired("input")

	serviceCmd.Flags().String("ff", "", "flag format, only prefix") // i.e if FLAG{...} then FLAG
	serviceCmd.MarkFlagRequired("ff")

	serviceCmd.Flags().String("zip", "", "output name of zip file, doesn't zip if not provided")

	/* useful for lazy people */
	serviceCmd.Flags().Bool("batch", false, "take defaults for all interactions")
	serviceCmd.Flags().Bool("cleanup", false, "remove temp folder after zipping")

	/* useful for known edge cases */
	serviceCmd.Flags().String("exclude", "", "exclude from flag check") // i.e FLAG{THIS_IS_A_FAKE_FLAG}
}
