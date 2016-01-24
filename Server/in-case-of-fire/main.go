// Author: Huy Pham
// About:  Very simple server for funny project: In Case Of Fire.

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	username = "NghiaTran"
)

func main() {
	h := Handle{}
	router := gin.Default()

	// For static and html file
	router.Static("/css", "css")
	router.LoadHTMLGlob("html/*")

	// Home page
	router.GET("/", h.Home)

	// Home page
	router.GET("/home", h.Home)

	// Remove page
	router.POST("/remove", h.Remove)

	// Remove page
	router.GET("/remove", h.Home)

	// Register project
	router.POST("/home", h.RegisterProject)

	// get project
	router.GET("/projects", h.GetProject)

	// In case of fire
	router.GET("/in-case-of-fire", h.InCaseOfFire)

	router.Run(":10000")
}

type Handle struct{}

func (h *Handle) HTML(c *gin.Context) {
	paths := strings.Split(c.Request.URL.Path, "/")
	c.HTML(http.StatusOK, paths[len(paths)-1]+".html", nil)
}

func (h *Handle) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func (h *Handle) Remove(c *gin.Context) {
	type Path struct {
		Path string `form:"path"`
	}
	path := Path{}
	c.Bind(&path)
	removeProjects(path.Path, c)
	c.HTML(http.StatusOK, "home.html", nil)
}

func (h *Handle) InCaseOfFire(c *gin.Context) {
	go SaveYourLife(username)
	c.String(200, "Okie")
}

func (h *Handle) RegisterProject(c *gin.Context) {
	type Path struct {
		Path string `form:"path"`
	}
	path := Path{}
	c.Bind(&path)
	paths := appendProjects(path.Path, c)
	CreateSaveScript(username, getIP(c), paths)
	c.HTML(http.StatusOK, "home.html", nil)
}

func (h *Handle) GetProject(c *gin.Context) {
	content, _ := ioutil.ReadFile("projects")
	if string(content) == "" {
		c.JSON(200, gin.H{"paths": []string{}})
		return
	}
	paths := strings.Split(string(content), ",")
	c.JSON(200, gin.H{"paths": paths})
}

func getIP(c *gin.Context) string {
	iPArray := strings.Split(c.ClientIP(), ":")
	pureIp := iPArray[0]
	return pureIp
}

// Helper function
func appendProjects(path string, c *gin.Context) []string {
	filename := "projects"
	content, _ := ioutil.ReadFile(filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	paths := []string{}
	if string(content) != "" {
		paths = strings.Split(string(content), ",")
	}

	isDup := false
	for _, p := range paths {
		if p == path || path == "" {
			isDup = true
		}
	}
	newPaths := []string{}
	if isDup {
		newPaths = paths
	} else {
		newPaths = append(paths, path)
	}
	n, err := io.WriteString(f, strings.Join(newPaths, ","))
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
	CreateSaveScript(username, getIP(c), newPaths)
	return newPaths
}

func removeProjects(removePath string, c *gin.Context) []string {
	filename := "projects"
	content, _ := ioutil.ReadFile(filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	paths := strings.Split(string(content), ",")
	newPaths := []string{}
	for _, path := range paths {
		fmt.Print(path)
		if path != removePath {
			newPaths = append(newPaths, path)
		}
	}
	n, err := io.WriteString(f, strings.Join(newPaths, ","))
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
	CreateSaveScript(username, getIP(c), newPaths)
	return newPaths
}

// Helper function
func SaveYourLife(name string) {
	fmt.Println("Exectute task:", name)
	_, err := exec.Command("./save_" + name).Output()
	if err == nil {
		fmt.Println("Save:", name, " success!")
	} else {
		fmt.Print(err)
	}
}

// Create bash script for pushing project
func CreateSaveScript(username, host string, paths []string) {
	content := "#!/bin/bash\nssh -t -t"
	content = content + " " + username + "@" + host + " << EOF\n"
	for _, path := range paths {
		if path == "" {
			continue
		}
		content = content + "	cd " + path + " && \\\n"
		content = content + "	git checkout -B fire && \\\n	git push origin fire\n"
	}
	content = content + "	exit\nEOF"
	filename := "save_" + username
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	n, err := io.WriteString(f, content)
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
	_, _ = exec.Command("sh", "-c", "chmod +x "+filename).Output()
}
