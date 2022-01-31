// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/web"
)

type StatusResponse struct {
	Status string `json:"Status"`
	Link   string `json:"Link"`
}

type Result struct {
	Output string
	Err    error
}

type Job struct {
	Id            int
	ResultChannel chan Result
}

type StatusDesc struct {
	Success bool           `json:"success"`
	Message StatusResponse `json:"message"`
	Errors  string         `json:"errors"`
}

var jobMap = make(map[int]Job)
var resultMap = make(map[int]Result)
var jobCounter int = 0

func NewJob() *Job {
	var job Job
	job.ResultChannel = make(chan Result)

	jobCounter++
	job.Id = jobCounter
	jobMap[jobCounter] = job

	return &job
}

func RemoveJob(id int) {
	delete(jobMap, id)
}

func RemoveResult(id int) {
	delete(resultMap, id)
}

func CreateJob(acquireFunc func() (string, error)) *Job {
	job := NewJob()
	go func() {
		s, err := acquireFunc()
		var result Result
		result.Output = s
		result.Err = err
		job.ResultChannel <- result
	}()
	return job
}

func AcceptedResponse(w http.ResponseWriter, job *Job) error {
	w.Header().Set("Location", "/api/v1/_jobs/status/"+strconv.Itoa(job.Id))
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func routerAcquireStatus(w http.ResponseWriter, r *http.Request) {
	var err error

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if job, ok := jobMap[id]; ok {
		select {
		case result := <-job.ResultChannel:
			resultMap[id] = result
			RemoveJob(id)
			web.JSONResponse(
				StatusResponse{
					Status: "complete",
					Link:   "/api/v1/_jobs/result/" + strconv.Itoa(id),
				},
				w)
		default:
			web.JSONResponse(StatusResponse{Status: "inprogress"}, w)
		}
	} else {
		err = errors.New("not found")
		web.JSONResponseError(err, w)
	}
}

func routerAcquireResult(w http.ResponseWriter, r *http.Request) {
	var err error

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if result, ok := resultMap[id]; ok {
		if result.Err != nil {
			web.JSONResponseError(result.Err, w)
		} else {
			if result.Output != "" {
				var r interface{}
				err := json.Unmarshal([]byte(result.Output), &r)
				if err != nil {
					RemoveResult(id)
					web.JSONResponseError(err, w)
				}
				web.JSONResponse(r, w)
			} else {
				web.JSONResponse("", w)
			}
		}
		RemoveResult(id)
	} else {
		err = errors.New("not found")
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterJobs(router *mux.Router) {
	n := router.PathPrefix("/_jobs").Subrouter().StrictSlash(false)

	n.HandleFunc("/status/{id}", routerAcquireStatus).Methods("GET")
	n.HandleFunc("/result/{id}", routerAcquireResult).Methods("GET")
}

func DispatchAndWait(method, host string, url string, token map[string]string, data interface{}) ([]byte, error) {
	var msg []byte
	r, err := web.DispatchSocketWithStatus(method, host, url, token, data)
	if err != nil {
		return nil, err
	}
	if r.StatusCode == 202 {
		if location := r.Header.Get("Location"); location != "" {
			for {
				s, err := web.DispatchSocket(http.MethodGet, host, location, token, nil)
				if err != nil {
					fmt.Printf("retrieving job status failed: %v\n", err)
					return nil, err
				}
				status := StatusDesc{}
				err = json.Unmarshal(s, &status)
				if err != nil {
					fmt.Printf("Failed to decode json message: %v\n", err)
					return nil, err
				}
				if status.Message.Status == "complete" {
					link := status.Message.Link
					msg, err = web.DispatchSocket(http.MethodGet, host, link, token, nil)
					if err != nil {
						fmt.Printf("retrieving result failed: %v\n", err)
						return nil, err
					}
					break
				} else if status.Message.Status != "inprogress" {
					err = errors.New("unexptected status")
					return nil, err
				}
				time.Sleep(1 * time.Second)
				fmt.Printf(".")
			}
			fmt.Printf("\n")
		} else {
			err = errors.New("no location in headers")
			return nil, err
		}
	} else if r.StatusCode == 200 {
		msg = r.Body
	} else {
		return nil, err
	}
	return msg, err
}
