package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetDefaultFiles(c *gin.Context) {
	var importantDirs []string
	var availableDrives []string
	dirNames := []string{
		"Downloads",
		"Documents",
		"Desktop",
		"Pictures",
		"Videos",
		"Music",
		"Applications",
		"Library",
	}
	type DirInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get user home directory"})
		return
	}

	switch runtime.GOOS {
	case "windows":
		importantDirs = []string{
			filepath.Join(homeDir, "Downloads"),
			filepath.Join(homeDir, "Documents"),
			filepath.Join(homeDir, "Desktop"),
			filepath.Join(homeDir, "Pictures"),
			filepath.Join(homeDir, "Videos"),
			filepath.Join(homeDir, "Music"),
		}
		availableDrives = getWindowsDrives()
	case "darwin": // macOS
		importantDirs = []string{
			filepath.Join(homeDir, "Downloads"),
			filepath.Join(homeDir, "Documents"),
			filepath.Join(homeDir, "Desktop"),
			filepath.Join(homeDir, "Pictures"),
			filepath.Join(homeDir, "Movies"),
			filepath.Join(homeDir, "Music"),
			filepath.Join(homeDir, "Library"),
			"/Applications",
		}
		availableDrives = getUnixMountPoints()
	case "linux":
		importantDirs = []string{
			filepath.Join(homeDir, "Downloads"),
			filepath.Join(homeDir, "Documents"),
			filepath.Join(homeDir, "Desktop"),
			filepath.Join(homeDir, "Pictures"),
			filepath.Join(homeDir, "Videos"),
			filepath.Join(homeDir, "Music"),
		}
		availableDrives = getUnixMountPoints()
	default:
		c.JSON(500, gin.H{"error": "Unsupported operating system"})
		return
	}

	// Filter out non-existent directories

	existingDirs := make([]DirInfo, 0)
	for _, dir := range importantDirs {
		for _, name := range dirNames {
			if strings.Contains(dir, name) {
				if _, err := os.Stat(dir); err == nil {
					existingDirs = append(existingDirs, DirInfo{Name: name, Path: dir})
				} else {
					existingDirs = append(existingDirs, DirInfo{Name: name, Path: ""})
				}
				break
			}
		}
	}

	c.JSON(200, gin.H{
		"os":          runtime.GOOS,
		"directories": existingDirs,
		"drives":      availableDrives,
	})
}

func getWindowsDrives() []string {
	drives := []string{}
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			drives = append(drives, string(drive)+":\\\\")
			f.Close()
		}
	}
	return drives
}

func getUnixMountPoints() []string {
	mountPoints := []string{}
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return mountPoints
	}
	defer file.Close()

	var device, mountPoint, fsType, options string
	var dump, pass int
	for {
		_, err := fmt.Fscanf(file, "%s %s %s %s %d %d\n", &device, &mountPoint, &fsType, &options, &dump, &pass)
		if err != nil {
			break
		}
		if device != "none" && !filepath.HasPrefix(mountPoint, "/proc") && !filepath.HasPrefix(mountPoint, "/sys") {
			mountPoints = append(mountPoints, mountPoint)
		}
	}
	return mountPoints
}
