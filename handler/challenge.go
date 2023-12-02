package handler

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type challengeHandler struct {
	challengeService entity.ChallengeService
}

func NewChallengeHandler(r *gin.RouterGroup, service entity.ChallengeService) {
	handler := &challengeHandler{
		challengeService: service,
	}

	r.GET("/challenge/labs", handler.GetAllLabsRedacted)
	r.GET("/challenge/labs/my", handler.GetMyChallengeLabsRedacted)
	r.GET("/challenge", handler.GetAllChallenges)
	r.GET("/challenge/my", handler.GetMyChallenges)
	r.POST("/challenge", handler.UpsertChallenge)
	r.DELETE("/challenge/:challengeId", handler.DeleteChallenge)
}

func (ch *challengeHandler) GetAllLabsRedacted(c *gin.Context) {
	labs, err := ch.challengeService.GetAllLabsRedacted()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, labs)
}

func (ch *challengeHandler) GetMyChallengeLabsRedacted(c *gin.Context) {

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	labs, err := ch.challengeService.GetChallengesByUserId(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, labs)
}

func (ch *challengeHandler) GetAllChallenges(c *gin.Context) {
	challenges, err := ch.challengeService.GetAllChallenges()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, challenges)
}

func (ch *challengeHandler) GetMyChallenges(c *gin.Context) {

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	challenges, err := ch.challengeService.GetChallengesByUserId(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, challenges)
}

func (ch *challengeHandler) UpsertChallenge(c *gin.Context) {
	challenge := entity.Challenge{}
	if err := c.BindJSON(&challenge); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := ch.challengeService.UpsertChallenge(challenge); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, challenge)
}

func (ch *challengeHandler) DeleteChallenge(c *gin.Context) {
	challengeId := c.Param("challengeId")

	challengeIds := []string{challengeId}

	if err := ch.challengeService.DeleteChallenges(challengeIds); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, challengeId)
}
