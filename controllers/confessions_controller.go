package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gosu-team/fptu-api/lib"
	"github.com/gosu-team/fptu-api/models"
	recaptcha "gopkg.in/ezzarghili/recaptcha-go.v3"
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

type confessionElement struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at, omitempty"`
	UpdatedAt *time.Time `json:"updated_at, omitempty"`
	Content   string     `json:"content"`
	Status    int        `json:"status"`
	Approver  string     `json:"approver"`
	Reason    string     `json:"reason"`
	CfsID     int        `json:"cfs_id"`
}

type confessionCollection []confessionElement

func confessionsResponseResolve(arr []models.Confession) confessionCollection {
	var collection confessionCollection
	for _, e := range arr {
		// Mapping approver email
		user := new(models.User)
		approverNickname := user.FetchNicknameByID(e.Approver)

		e := confessionElement{e.ID, e.CreatedAt, e.UpdatedAt, e.Content, e.Status, approverNickname, e.Reason, e.CfsID}
		collection = append(collection, e)
	}

	return collection
}

// GetAllConfessionsHandler ...
func GetAllConfessionsHandler(w http.ResponseWriter, r *http.Request) {
	// Number of element to query
	numLoad, err := strconv.Atoi(r.URL.Query().Get("load"))
	if err != nil {
		numLoad = 10
	}

	res := lib.Response{ResponseWriter: w}
	confession := new(models.Confession)
	confessions := confession.FetchAll(numLoad)

	collection := confessionsResponseResolve(confessions)

	res.SendOK(collection)
}

type requestConfessionBySenderRequest struct {
	Token string `json:"token"`
}

// GetConfessionsBySenderHandler ...
func GetConfessionsBySenderHandler(w http.ResponseWriter, r *http.Request) {
	// Number of element to query
	numLoad, err := strconv.Atoi(r.URL.Query().Get("load"))
	if err != nil {
		numLoad = 10
	}

	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	confession := new(models.Confession)
	tokenRequest := new(requestConfessionBySenderRequest)
	req.GetJSONBody(tokenRequest)

	confessions := confession.FetchBySender(tokenRequest.Token, numLoad)
	collection := confessionsResponseResolve(confessions)
	res.SendOK(collection)
}

type newConfessionRequest struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
	Captcha string `json:"captcha"`
}

// CreateConfessionHandler ...
func CreateConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	newConfession := new(newConfessionRequest)
	req.GetJSONBody(newConfession)

	// Verify captcha
	recaptchaSecret := os.Getenv("CAPTCHA")
	captcha, _ := recaptcha.NewReCAPTCHA(recaptchaSecret, recaptcha.V2, 10*time.Second)

	err := captcha.Verify(newConfession.Captcha)
	if err != nil {
		// Track log
		fmt.Println("Captcha secrect: ", recaptchaSecret)
		fmt.Println("Captcha string: ", newConfession.Captcha)

		res.SendBadRequest("Invalid captcha!")
		return
	}

	if len(newConfession.Content) < 1 {
		res.SendBadRequest("Too short!")
		return
	}

	confession := models.Confession{
		Content: newConfession.Content,
		Sender:  newConfession.Sender,
	}

	if err := confession.Create(); err != nil {
		res.SendBadRequest(err.Error())
		return
	}

	res.SendCreated(confession)
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

// RollbackApproveConfessionHandler ...
func RollbackApproveConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	approverID := getUserIDFromHeader(r)
	approveConfession := new(models.Confession)
	req.GetJSONBody(approveConfession)

	if err := approveConfession.RollbackApproveConfession(approverID); err != nil {
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

// NextConfession ...
type NextConfession struct {
	ID int `json:"id"`
}

// GetNextConfessionNextIDHandler ...
func GetNextConfessionNextIDHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}

	confession := new(models.Confession)
	nextID := confession.GetNextConfessionID()
	nextConfession := NextConfession{
		ID: nextID,
	}

	res.SendOK(nextConfession)
}
