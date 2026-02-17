package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
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
	Size        string `json:"size,omitempty"`
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
			_, perms, err := getPerms(d)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := ActionResponse{
			CommandOutput: cmdLog.String(),
		}
		c.JSON(http.StatusOK, response)
	})

	router.POST("/api/action/7zzip001", func(c *gin.Context) {
		output, err := run7z(downloadsDirectory, "-y", "x", "*.zip.001")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": output + " - " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, ActionResponse{output})
	})

	router.POST("/api/action/7zzip", func(c *gin.Context) {
		c.String(http.StatusOK, "7z -y x \"*.zip\" output")
	})

	router.POST("/api/action/7z7z001", func(c *gin.Context) {
		c.String(http.StatusOK, "7z -y x \"*.7z.001\" output")
	})

	router.POST("/api/action/rmzip", func(c *gin.Context) {
		cmdLog := strings.Builder{}
		matches, err := filepath.Glob(downloadsDirectory + "/*zip*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, match := range matches {
			_, err := fmt.Fprintf(&cmdLog, "rm %s\n", match)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			err = os.Remove(match)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		response := ActionResponse{
			CommandOutput: cmdLog.String(),
		}
		c.JSON(http.StatusOK, response)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fileTree := make([]FileTreeItem, 0)
		for _, entry := range dir {
			size, perms, err := getPerms(entry)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			hrSize := humanize.Bytes(uint64(size))
			fileTree = append(fileTree, FileTreeItem{
				Name:        entry.Name(),
				Size:        hrSize,
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

func getPerms(f os.DirEntry) (int64, os.FileMode, error) {
	info, err := f.Info()
	if err != nil {
		return 0, 0, err
	}
	return info.Size(), info.Mode().Perm(), nil
}

func run7z(directory string, args ...string) (string, error) {
	cmd := exec.Command("/bin/7z", args...)
	cmd.Dir = directory
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	output, err := cmd.Output()
	if err != nil {
		slurp, _ := io.ReadAll(stderr)
		return cmd.String() + ": " + string(slurp), err
	}
	return cmd.String() + "\n" + string(output), err
}
