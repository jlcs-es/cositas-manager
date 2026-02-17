package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"net/http"
)

//go:embed cositas-manager/dist/cositas-manager/browser
var frontendFS embed.FS

type MoveFileBody struct {
	FileName string `json:"fileName" binding:"required"`
	DirName  string `json:"dirName" binding:"required"`
}

type ActionResponse struct {
	CommandOutput string `json:"commandOutput"`
}

type FileTreeItem struct {
	Name        string `json:"name,omitempty"`
	Permissions string `json:"permissions,omitempty"`
	IsDirectory bool   `json:"isDirectory,omitempty"`
}

func main() {

	downloadsDirectory, ok := os.LookupEnv("DOWNLOADS_DIRECTORY")
	if !ok {
		downloadsDirectory = "/var/home/joseluis/Descargas"
	}

	router := gin.Default()

	// Serve embedded frontend files

	sub, err := fs.Sub(frontendFS, "cositas-manager/dist/cositas-manager/browser")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/", http.FS(sub))

	// Backend API logic

	router.POST("/api/action/chmod", func(c *gin.Context) {
		cmdLog := strings.Builder{}
		err := filepath.WalkDir(downloadsDirectory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			perms, err := getPerms(d)
			if err != nil {
				return err
			}
			noexec := perms & 0666
			_, err = fmt.Fprintf(&cmdLog, "chmod %s -> %s : %s\n", perms, noexec, d.Name())
			if err != nil {
				return err
			}
			err = os.Chmod(path, noexec)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		response := ActionResponse{
			CommandOutput: cmdLog.String(),
		}
		c.JSON(http.StatusOK, response)
	})

	router.POST("/api/action/7zzip001", func(c *gin.Context) {
		c.String(http.StatusOK, "7z -y x \"*.zip.001\" output")
	})

	router.POST("/api/action/7zzip", func(c *gin.Context) {
		c.String(http.StatusOK, "7z -y x \"*.zip\" output")
	})

	router.POST("/api/action/7z7z001", func(c *gin.Context) {
		c.String(http.StatusOK, "7z -y x \"*.7z.001\" output")
	})

	router.POST("/api/action/rmzip", func(c *gin.Context) {
		c.String(http.StatusOK, "rm *zip* output")
	})

	router.POST("/api/action/move", func(c *gin.Context) {
		body := MoveFileBody{}
		if err := c.ShouldBind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		response := ActionResponse{
			CommandOutput: "move " + body.FileName + " to " + body.DirName,
		}
		c.JSON(http.StatusOK, response)
	})

	router.POST("/api/info/listfiles", func(c *gin.Context) {
		dir, err := os.ReadDir(downloadsDirectory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		fileTree := make([]FileTreeItem, 0)
		for _, entry := range dir {
			perms, err := getPerms(entry)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			fileTree = append(fileTree, FileTreeItem{
				Name:        entry.Name(),
				Permissions: perms.String(),
				IsDirectory: entry.IsDir(),
			})
		}
		c.JSON(http.StatusOK, fileTree)
	})

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getPerms(f os.DirEntry) (os.FileMode, error) {
	info, err := f.Info()
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}
