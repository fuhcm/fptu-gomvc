package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gosu-team/cfapp-api/lib"
	"github.com/gosu-team/cfapp-api/models"
)

func getUserIDFromHeader(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	tokenString := bearerToken[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0
	}

	if claims, valid := token.Claims.(jwt.MapClaims); valid && token.Valid {
		userIDStr := claims["jti"].(string)
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			userID = 0
		}

		return userID
	}

	return 0
}

// GetAllConfessionsHandler ...
func GetAllConfessionsHandler(w http.ResponseWriter, r *http.Request) {
	// Number of element to query
	numLoad, err := strconv.Atoi(r.URL.Query().Get("numLoad"))
	if err != nil {
		numLoad = 10
	}

	res := lib.Response{ResponseWriter: w}
	confession := new(models.Confession)
	confessions := confession.FetchAll(numLoad)
	res.SendOK(confessions)
}

// CreateConfessionHandler ...
func CreateConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	confession := new(models.Confession)
	req.GetJSONBody(confession)

	if err := confession.Create(); err != nil {
		res.SendBadRequest(err.Error())
		return
	}

	res.SendCreated(confession)
}

// GetConfessionsBySenderHandler ...
func GetConfessionsBySenderHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	confession := new(models.Confession)
	req.GetJSONBody(confession)
	confessions := confession.FetchBySender(confession.Sender)
	res.SendOK(confessions)
}

// Overview ...
type Overview struct {
	TotalConfess    int `json:"total"`
	PendingConfess  int `json:"pending"`
	RejectedConfess int `json:"rejected"`
}

// GetConfessionsOverviewHandler ...
func GetConfessionsOverviewHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}

	confession := new(models.Confession)
	totalConfess, pendingConfess, rejectedConfess := confession.FetchOverview()
	overview := Overview{
		TotalConfess:    totalConfess,
		PendingConfess:  pendingConfess,
		RejectedConfess: rejectedConfess,
	}

	res.SendOK(overview)
}

// ApproveConfessionHandler ...
func ApproveConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	approverID := getUserIDFromHeader(r)
	approveConfession := new(models.Confession)
	req.GetJSONBody(approveConfession)

	if err := approveConfession.ApproveConfession(approverID); err != nil {
		res.SendBadRequest(err.Error())
		return
	}

	res.SendOK(approveConfession)
}

// RejectConfessionRequest ...
type RejectConfessionRequest struct {
	ID     int    `json: "id"`
	Reason string `json: "reason"`
}

// RejectConfessionHandler ...
func RejectConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	approverID := getUserIDFromHeader(r)
	rejectConfessionRequest := new(RejectConfessionRequest)
	req.GetJSONBody(rejectConfessionRequest)

	rejectConfession := models.Confession{ID: rejectConfessionRequest.ID}

	if err := rejectConfession.RejectConfession(approverID, rejectConfessionRequest.Reason); err != nil {
		res.SendBadRequest(err.Error())
		return
	}

	res.SendOK(rejectConfession)
}
