package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"webserver/lib"
	"webserver/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	recaptcha "gopkg.in/ezzarghili/recaptcha-go.v2"
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
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Content   string     `json:"content"`
	Status    int        `json:"status"`
	Approver  string     `json:"approver"`
	Reason    string     `json:"reason"`
	CfsID     int        `json:"cfs_id"`
	PushID    string     `json:"push_id"`
}

type confessionCollection []confessionElement

func confessionsResponseResolve(arr []models.Confession) confessionCollection {
	var collection confessionCollection
	for _, e := range arr {
		// Mapping approver email
		user := new(models.User)
		approverNickname := user.FetchNicknameByID(e.Approver)

		e := confessionElement{e.ID, e.CreatedAt, e.UpdatedAt, e.Content, e.Status, approverNickname, e.Reason, e.CfsID, e.PushID}
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

// GetApprovedConfessionsHandler ...
func GetApprovedConfessionsHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}

	// Number of element to query
	latestID, err := strconv.Atoi(r.URL.Query().Get("latest_id"))
	if err != nil {
		latestID = 0
	}

	confession := new(models.Confession)
	confessions := confession.FetchApprovedConfession(latestID)
	collection := confessionsResponseResolve(confessions)
	res.SendOK(collection)
}

type newConfessionRequest struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
	Captcha string `json:"captcha"`
	PushID  string `json:"pushid"`
}

// CreateConfessionHandler ...
func CreateConfessionHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	newConfession := new(newConfessionRequest)
	req.GetJSONBody(newConfession)

	// Verify captcha
	recaptchaSecret := os.Getenv("CAPTCHA")
	captcha, _ := recaptcha.NewReCAPTCHA(recaptchaSecret)

	bool, err := captcha.VerifyNoRemoteIP(newConfession.Captcha)

	if bool != true || err != nil {
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
		PushID:  newConfession.PushID,
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

// RejectConfessionRequest ...
type RejectConfessionRequest struct {
	ID     int    `json:"id"`
	Reason string `json:"reason"`
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

// SearchConfessionsHandler ...
func SearchConfessionsHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}

	// Number of element to query
	keyword := r.URL.Query().Get("q")
	keyword = strings.TrimSpace(keyword)

	if keyword == "" {
		res.SendBadRequest("Nothing to search!")
		return
	}

	confession := new(models.Confession)
	confessions := confession.SearchConfession(keyword)
	collection := confessionsResponseResolve(confessions)
	res.SendOK(collection)
}

// SyncPushIDRequest ...
type SyncPushIDRequest struct {
	Sender string `json:"sender"`
	PushID string `json:"push_id"`
}

// SyncPushIDHandler ...
func SyncPushIDHandler(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	jsonRequest := new(SyncPushIDRequest)
	req.GetJSONBody(jsonRequest)

	confession := new(models.Confession)
	confession.SyncPushID(jsonRequest.Sender, jsonRequest.PushID)

	res.SendNoContent()
}

// RadioType ...
type RadioType struct {
	Radios string `json:"radios"`
}

// GetRedisClient ...
func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
}

// RedisRead ...
func RedisRead(key string) (string, error) {
	client := GetRedisClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

// RedisWrite ...
func RedisWrite(key string, value string) error {
	client := GetRedisClient()

	err := client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// SetRadio ...
func SetRadio(w http.ResponseWriter, r *http.Request) {
	req := lib.Request{ResponseWriter: w, Request: r}
	res := lib.Response{ResponseWriter: w}

	radioRequest := new(RadioType)
	req.GetJSONBody(radioRequest)

	err := RedisWrite("radios", radioRequest.Radios)
	if err != nil {
		logrus.Println(err.Error())
	}

	res.SendOK(radioRequest)
}

// GetRadio ...
func GetRadio(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}

	value, err := RedisRead("radios")
	if err != nil {
		logrus.Println(err.Error())
	}

	radiosObj := RadioType{Radios: value}

	res.SendOK(radiosObj)
}
