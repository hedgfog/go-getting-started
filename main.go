package main

import (
	"context"
	"github.com/heroku/go-getting-started/servcies/gitsource/githubsource"
	"golang.org/x/oauth2"
	oauthgit "golang.org/x/oauth2/github"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

var (
	gitAuth = &oauth2.Config{
		ClientID:     "35ccfc4585abf2e301c7",
		ClientSecret: "aba1d30229dce48152d2c4cc309f7f09020b4b68",
		RedirectURL:  "http://localhost:8080/account/github/callback",
		Scopes: []string{
			"user:email",
		},
		Endpoint: oauthgit.Endpoint,
	}
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/link/github", ghInit)
	router.GET("/account/github/callback", ghCallback)
	router.GET("/account/github/workbook", handleWebhook)

	router.Run(":" + port)
}

func ghInit(c *gin.Context) {
	url := gitAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

func ghCallback(c *gin.Context) {
	code := c.Query("code")
	tok, err := gitAuth.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	client := githubsource.New(gitAuth.Client(context.Background(), tok))

	info, err := client.ListUserRepos()

	err = client.CreateRepoWebhook("hedgfog/relayer", "http://localhost:8080//webhook", "SECRET")
	if err != nil {
		log.Println("ERR:", err)
	}

	c.JSON(200, info)
}

func handleWebhook(c *gin.Context) {

}
