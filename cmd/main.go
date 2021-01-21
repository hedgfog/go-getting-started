package main

import (
	"cilintservice/servcies/gitsource/githubsource"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	oauthgit "golang.org/x/oauth2/github"
	"log"
	"os"
)

var (
	gitAuth = &oauth2.Config{
		ClientID:     "35ccfc4585abf2e301c7",
		ClientSecret: "aba1d30229dce48152d2c4cc309f7f09020b4b68",
		RedirectURL:  "https://cilinter.herokuapp.com/account/github/callback",
		Scopes: []string{
			"user:email",
			"write:repo_hook",
			"repo",
			"admin:repo_hook",
			"public_repo",
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
	router.Use(gin.Recovery())

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

	err = client.CreateRepoWebhook("hedgfog/relayer", "https://cilinter.herokuapp.com/account/github/workbook", "SECRET")
	if err != nil {
		log.Println("ERR:", err)
	}

	c.JSON(200, info)
}

func handleWebhook(c *gin.Context) {

}
