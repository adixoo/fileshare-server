package helpers

import (
	"log"
	"os"
)

func ReadDirectory(root string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var files []map[string]interface{}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			continue
		}

		fileType := ""
		fileSize := int64(0)
		isDirectory := entry.IsDir()
		modTime := info.ModTime()

		if !isDirectory {
			fileType = info.Name()[len(info.Name())-3:]
			fileSize = info.Size()
		}

		files = append(files, map[string]interface{}{
			"name":         info.Name(),
			"directory":    isDirectory,
			"type":         fileType,
			"size":         fileSize / 1024, // Convert fileSize to KB
			"modification": modTime,
		})
	}

	return files, nil
}
