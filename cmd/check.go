package cmd

import (
	"fmt"
	"strings"

	"github.com/gatariee/sanity/internal/check"
	"github.com/gatariee/sanity/internal/logging"
	"github.com/spf13/cobra"

	
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check whether you're going to get fired",
	Long:  `i sure wonder whether i left any flags here`,
	Run: func(cmd *cobra.Command, args []string) {
		zip, _ := cmd.Flags().GetString("zip")
		dir, _ := cmd.Flags().GetString("dir")
		file, _ := cmd.Flags().GetString("file")

		if zip == "" && dir == "" && file == "" {
			fmt.Println("provide a zip, dir, or file")
			return
		}

		ff, _ := cmd.Flags().GetString("flag_format")

		if !strings.ContainsRune(ff, '{') {
			logging.LogWarn("flag format does not contain '{', adding it automatically.")
			ff = ff + "{"
		}

		if zip != "" {
			err := check.CheckZip(zip, ff)
			if err != nil {
				panic(err)
			}
		}

		if dir != "" {
			err := check.CheckDir(dir, ff)
			if err != nil {
				panic(err)
			}
		}

		if file != "" {
			err := check.CheckFile(file, ff)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	/* optional flags */
	checkCmd.Flags().String("zip", "", "zip file to check (i.e dist.zip)")
	checkCmd.Flags().String("dir", "", "directory to check (i.e ./dist)")
	checkCmd.Flags().String("file", "", "file to check (i.e ./dist/index.html)")

	/* required flags */
	checkCmd.Flags().String("flag_format", "", "flag format to check (i.e FLAG{...} -> FLAG)")
	checkCmd.MarkFlagRequired("flag_format")
}
