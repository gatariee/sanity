package check

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gatariee/sanity/internal/logging"
)

func CheckDir(dir string, ff string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		CheckFile(path, ff)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckZip(zip string, ff string) error {
	fmt.Println("checking zip")
	// not implemented, lazy
	return nil
}

func CheckFile(file string, ff string) error {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	fc, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	if len(fc) == 0 {
		return fmt.Errorf("file is empty")
	}

	ctn := string(fc)
	idx := strings.Index(ctn, ff)
	if idx == -1 {
		return nil
	}

	logging.LogWarn("flag found in file %s, on line %d", absPath, strings.Count(ctn[:idx], "\n")+1)

	sp := idx + len(ff)
	logging.LogInfo("flag (+10 chars): %s%s", ff, ctn[sp:sp+10])

	return nil
}
