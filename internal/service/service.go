package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gatariee/sanity/internal/logging"
)

func CheckFlagFile(absPath string, flagFormat string, exclude string) error {
	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "flag.txt" || info.Name() == "flag" {
			/**
			We found a "flag.txt" or a "flag" file.
			*/
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if len(content) == 0 {
				return fmt.Errorf("flag file is empty")
			}

			fakeFlag := fmt.Sprintf("%s{THIS_IS_A_FAKE_FLAG}", flagFormat)

			logging.LogWarn("found potential flag file at: %s", path)
			logging.LogInfo("content: %s", string(content))

			if exclude != "" {
				if bytes.Contains(content, []byte(exclude)) {
					logging.LogWarn("excluded flag found in file, skipping")
					return nil
				}
			}

			logging.LogInfo("replace the flag with fake flag \"%s\"? (y/n)", fakeFlag)

			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				logging.LogInfo("ok, replacing flag file with fake flag")
				err := os.WriteFile(path, []byte(fakeFlag), 0o644)
				if err != nil {
					return err
				}
			} else {
				logging.LogWarn("leaving flag file as is, be careful")

				/**
				we can afford to return early here cuz we checked the contents of the file already
				*/
				return nil

			}
		}

		if info.Size() > 0 {
			fileContent := make([]byte, info.Size())
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = io.ReadFull(file, fileContent)
			if err != nil {
				return err
			}

			/**
			Did we accidentally leave a flag in the file?
			*/
			if bytes.Contains(fileContent, []byte(flagFormat)) {
				logging.LogWarn("found potential flag in file at: %s", path)
				logging.LogInfo("content: %s", string(fileContent))

				if exclude != "" {
					if bytes.Contains(fileContent, []byte(exclude)) {
						logging.LogWarn("excluded flag found in file, skipping")
						return nil
					}
				}

				logging.LogInfo("replace the flag with fake flag \"%s\"? (y/n)", "THIS_IS_A_FAKE_FLAG")

				var response string
				fmt.Scanln(&response)

				if response == "y" || response == "Y" {
					logging.LogInfo("ok, replacing flag in file with fake flag")
					newContent := bytes.ReplaceAll(fileContent, []byte(flagFormat), []byte("THIS_IS_A_FAKE_FLAG"))
					err := os.WriteFile(path, newContent, os.ModePerm)
					if err != nil {
						return err
					}
				} else {
					logging.LogWarn("leaving flag in file as is, be careful")
					return nil
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func PrepareService(servicePath string, flagFormat string, exclude string) error {
	_, err := os.Stat(servicePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("service folder does not exist")
	}

	var absPath string
	if filepath.IsAbs(servicePath) {
		absPath = servicePath
	} else {
		absPath, err = filepath.Abs(servicePath)
		if err != nil {
			return err
		}
	}

	/**
	Did we accidentally leave a flag in the service folder, or in the contents of a file?
	*/
	err = CheckFlagFile(absPath, flagFormat, exclude)
	if err != nil {
		return err
	}

	return nil
}
