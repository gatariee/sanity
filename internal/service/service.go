package service

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gatariee/sanity/internal/logging"
)

func CopyDir(dst, src string) error {
	/**
	https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
	*/
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil
		}

		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		in, _ := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		fh.Chmod(info.Mode())

		_, err = io.Copy(fh, in)
		return err
	})
}

func CheckFlagFile(absPath string, flagFormat string, exclude string, batch bool) error {
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

			fakeFlag := fmt.Sprintf("%sTHIS_IS_A_FAKE_FLAG}", flagFormat)

			logging.LogWarn("found potential flag file at: %s", path)
			logging.LogInfo("content: %s", string(content))

			if batch {
				logging.LogInfo("replacing flag file with fake flag \"%s\"", fakeFlag)
				logging.LogNewLine()
				err := os.WriteFile(path, []byte(fakeFlag), 0o644)
				if err != nil {
					return err
				}
				return nil
			}

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

		if info.Size() > 0 && !info.IsDir() {
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

				if batch {
					fakeFlag := fmt.Sprintf("%sTHIS_IS_A_FAKE_FLAG}", flagFormat)
					byteFakeFlag := []byte(fakeFlag)
					logging.LogInfo("replacing flag in file with fake flag \"%s\"", fakeFlag)
					logging.LogNewLine()
					err := os.WriteFile(path, byteFakeFlag, os.ModePerm)
					if err != nil {
						return err
					}
					return nil
				}

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
					fakeFlag := fmt.Sprintf("%sTHIS_IS_A_FAKE_FLAG}", flagFormat)
					err := os.WriteFile(path, []byte(fakeFlag), os.ModePerm)
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

func PrepareService(servicePath string, flagFormat string, exclude string, batch bool, tempPath string) error {
	
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

	if !strings.ContainsRune(flagFormat, '{') {
		logging.LogWarn("flag format does not contain '{', adding it automatically.")
		flagFormat = flagFormat + "{"
	}

	/**
	We don't want to mutate the service folder, so let's work directly on dist
	*/
	err = os.MkdirAll(tempPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = CopyDir(tempPath, absPath)
	if err != nil {
		return err
	}

	distAbsPath, err := filepath.Abs(tempPath)
	if err != nil {
		return err
	}

	/**
	Did we accidentally leave a flag in the service folder, or in the contents of a file?
	*/

	err = CheckFlagFile(distAbsPath, flagFormat, exclude, batch)
	if err != nil {
		return err
	}

	return nil
}
