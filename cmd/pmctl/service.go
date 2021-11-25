package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/systemd"
)

type UnitStatus struct {
	Success bool               `json:"success"`
	Message systemd.UnitStatus `json:"message"`
	Errors  string             `json:"errors"`
}

func executeSystemdUnitCommand(command string, unit string, host string, token map[string]string) {
	c := systemd.UnitAction{
		Action: command,
		Unit:   unit,
	}

	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/service/systemd", token, c)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from remote host: %v\n", err)
			return
		}
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/service/systemd", c)
		if err != nil {
			fmt.Printf("Failed to execute '%s': %v\n", command, err)
			os.Exit(1)
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to execute command: %v\n", m.Errors)
		os.Exit(1)
	}
}

func acquireSystemdUnitStatus(unit string, host string, token map[string]string) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/service/systemd/"+unit+"/status", token, nil)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from remote host: %v\n", err)
			return
		}
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/service/systemd/"+unit+"/status", nil)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from unix domain socket: %v\n", err)
			return
		}
	}

	u := UnitStatus{}
	if err := json.Unmarshal(resp, &u); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if u.Success {
		fmt.Printf("                  Name: %+v \n", u.Message.Name)
		fmt.Printf("           Description: %+v \n", u.Message.Description)
		fmt.Printf("               MainPid: %+v \n", u.Message.MainPid)
		fmt.Printf("             LoadState: %+v \n", u.Message.LoadState)
		fmt.Printf("           ActiveState: %+v \n", u.Message.ActiveState)
		fmt.Printf("              SubState: %+v \n", u.Message.SubState)
		fmt.Printf("         UnitFileState: %+v \n", u.Message.UnitFileState)

		if u.Message.StateChangeTimestamp > 0 {
			t := time.Unix(int64(u.Message.StateChangeTimestamp), 0)
			fmt.Printf("  StateChangeTimeStamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("  StateChangeTimeStamp: %+v \n", 0)
		}

		if u.Message.ActiveEnterTimestamp > 0 {
			t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
			fmt.Printf("  ActiveEnterTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("  ActiveEnterTimestamp: %+v \n", 0)
		}

		if u.Message.ActiveEnterTimestamp > 0 {
			t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
			fmt.Printf(" InactiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf(" InactiveExitTimestamp: %+v \n", 0)
		}

		if u.Message.ActiveExitTimestamp > 0 {
			t := time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			fmt.Printf("   ActiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("   ActiveExitTimestamp: %+v \n", 0)
		}

		if u.Message.InactiveExitTimestamp > 0 {
			t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
			fmt.Printf(" InactiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf(" InactiveExitTimestamp: %+v \n", 0)
		}

		switch u.Message.ActiveState {
		case "active", "reloading":
			if u.Message.ActiveEnterTimestamp > 0 {
				t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				fmt.Printf("                Active: %s (%s) since %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("                Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		case "inactive", "failed":
			if u.Message.ActiveExitTimestamp != 0 {
				t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
				fmt.Printf("                Active: %s (%s) since %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("                Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		case "activating":
			var t time.Time

			if u.Message.ActiveExitTimestamp > 0 || u.Message.ActiveEnterTimestamp > 0 {
				if u.Message.ActiveExitTimestamp > 0 {
					t = time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				} else if u.Message.ActiveEnterTimestamp > 0 {
					t = time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				}

				fmt.Printf("               Active: %s (%s) %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("               Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		default:
			t := time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			fmt.Printf("               Active: %s (%s) ago %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
		}
	} else {
		fmt.Println(u.Errors)
	}
}
