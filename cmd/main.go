package cmd

import (
	"cilintservice/servcies/gitsource/githubsource"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	oauthgit "golang.org/x/oauth2/github"
	"log"
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
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/link/github", ghInit)

	router.Run(":8080")

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
