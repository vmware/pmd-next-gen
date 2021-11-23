package login

import (
	"context"
	"net/http"

	"github.com/pm-web/pkg/web"
)

var loginMethod = map[string]string{
	"list-sessions":     "ListSessions",
	"list-users":        "ListUsers",
	"lock-session":      "LockSession",
	"lock-sessions":     "LockSessions",
	"terminate-session": "TerminateSession",
	"terminate-user":    "TerminateUser",
}

type Login struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

type User struct {
	UID  uint32 `json:"UID"`
	Name string `json:"Name"`
	Path string `json:"Path"`
}

type Session struct {
	ID   string `json:"ID"`
	UID  uint32 `json:"UID"`
	User string `json:"User"`
	Seat string `json:"Seat"`
	Path string `json:"Path"`
}

func AcquireUsersFromLogin(ctx context.Context, w http.ResponseWriter) error {
	users, err := DBusAcquireUsersFromLogin(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(users, w)
}

func AcquireSessionsFromLogin(ctx context.Context, w http.ResponseWriter) error {
	users, err := DBusAcquireUSessionsFromLogin(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(users, w)
}