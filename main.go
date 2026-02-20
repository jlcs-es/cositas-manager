package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
)

//go:embed cositas-manager/dist/cositas-manager/browser
var frontendFS embed.FS

type MoveActionBody struct {
	SourceName           string `json:"sourceName" binding:"required"`
	DestinationDirectory string `json:"destinationDirectory" binding:"required"`
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

type Job struct {
	output string
	date   time.Time
}

var jobsCounter atomic.Uint64
var jobs = sync.Map{}

func main() {

	go func() {
		for {
			jobs.Range(func(key, value any) bool {
				job := value.(*Job)
				if job.date.Before(time.Now().Add(-2 * time.Hour)) {
					jobs.Delete(key)
				}
				return true
			})
			time.Sleep(time.Hour)
		}
	}()

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
		jobID := run7zJob(downloadsDirectory, "-y", "x", "*.zip.001")
		msg := fmt.Sprintf("Working asynchronously - Job ID: %d", jobID)
		c.JSON(http.StatusOK, ActionResponse{msg})
	})

	router.POST("/api/action/7zzip", func(c *gin.Context) {
		jobID := run7zJob(downloadsDirectory, "-y", "x", "*.zip")
		msg := fmt.Sprintf("Working asynchronously - Job ID: %d", jobID)
		c.JSON(http.StatusOK, ActionResponse{msg})
	})

	router.POST("/api/action/7z7z001", func(c *gin.Context) {
		jobID := run7zJob(downloadsDirectory, "-y", "x", "*.7z.001")
		msg := fmt.Sprintf("Working asynchronously - Job ID: %d", jobID)
		c.JSON(http.StatusOK, ActionResponse{msg})
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
		body := MoveActionBody{}
		if err := c.ShouldBind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = os.MkdirAll(body.DestinationDirectory, 0775)
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := os.Rename(
			path.Join(downloadsDirectory, body.SourceName),
			path.Join(body.DestinationDirectory, body.SourceName),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := ActionResponse{
			CommandOutput: "Moved " + body.SourceName + " to " + body.DestinationDirectory,
		}
		c.JSON(http.StatusOK, response)
	})

	/*
		router.POST("/api/info/treemedia", func(c *gin.Context) {
			movies, err := listTree(downloadsDirectory)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, movies)
		})
	*/

	router.POST("/api/info/treemedia", func(c *gin.Context) {
		movies, err := listTree("/media/Movies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		series, err := listTree("/media/Series")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		anime, err := listTree("/media/Anime")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileTree := append(movies, append(series, anime...)...)

		c.JSON(http.StatusOK, fileTree)
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

	router.POST("/api/info/listjobs", func(c *gin.Context) {
		jobsList := make([]string, 0)
		jobs.Range(func(k, v any) bool {
			job := v.(*Job)
			jobDesc := fmt.Sprintf("%d - %s - %s", k.(uint64), job.date.UTC().Format(time.RFC3339), job.output)
			jobsList = append(jobsList, jobDesc)
			return true
		})
		c.JSON(http.StatusOK, jobsList)
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

func run7zJob(directory string, args ...string) uint64 {
	jobID := jobsCounter.Add(1)

	jobs.Store(jobID, &Job{
		output: "",
		date:   time.Now(),
	})

	go func() {
		j, _ := jobs.Load(jobID)
		job := j.(*Job)
		output, err := run7z(directory, args...)
		msg := "FINISHED: " + output
		if err != nil {
			msg = "ERROR: " + output + " - " + err.Error()
		}
		log.Println(msg)
		job.output = msg
	}()

	return jobID
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

func listTree(directory string) ([]string, error) {
	directories := make([]string, 0)
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	return directories, err
}
