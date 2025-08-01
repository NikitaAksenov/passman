package directory

import (
	"fmt"
	"os"
)

func CheckAndCreateDir(path string) error {
	// Check if path is empty
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		// If dir does not exits - create one
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed creating dir: %s", err.Error())
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("failed getting path info: %s", err.Error())
	}

	if !file.IsDir() {
		return fmt.Errorf("path is not dir")
	}

	return nil
}
