//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	apiDir := "api/v1"
	files, err := ioutil.ReadDir(apiDir)
	if err != nil {
		log.Fatal(err)
	}

	oldStr := `c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})`
	newStr := `response.FailWithMessage(c, constants.SYSTEM_ERROR)`

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") || file.Name() == "controller.go" {
			continue
		}

		filePath := filepath.Join(apiDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", filePath, err)
			continue
		}

		strContent := string(content)
		if strings.Contains(strContent, oldStr) {
			strContent = strings.ReplaceAll(strContent, oldStr, newStr)

			// Ensure import exists
			if !strings.Contains(strContent, `"kama_chat_server/pkg/util/response"`) {
				strContent = strings.Replace(strContent, `"net/http"`, `"net/http"`+"\n\t\"kama_chat_server/pkg/util/response\"", 1)
			}

			err = ioutil.WriteFile(filePath, []byte(strContent), os.ModePerm)
			if err != nil {
				log.Printf("Error writing %s: %v", filePath, err)
			} else {
				fmt.Printf("Refactored %s\n", filePath)
			}
		}
	}
}
