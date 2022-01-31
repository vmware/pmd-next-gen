// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package jobs

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
)

type Result struct {
	Output string
	Err    error
}

type Job struct {
	ResultChannel chan Result
	Id            int
}

type Jobs struct {
	jobMap     map[int]Job
	resultMap  map[int]Result
	jobCounter int
	Mutex      *sync.Mutex
}

var jobs *Jobs

func New() *Jobs {
	if jobs != nil {
		return jobs
	} else {
		jobs = &Jobs{
			jobMap:    make(map[int]Job),
			resultMap: make(map[int]Result),
			Mutex:     &sync.Mutex{},
		}
		return jobs
	}
}

func NewJob() *Job {
	jobs.Mutex.Lock()
	defer jobs.Mutex.Unlock()

	jobs.jobCounter++
	job := Job{
		ResultChannel: make(chan Result),
		Id:            jobs.jobCounter,
	}

	jobs.jobMap[jobs.jobCounter] = job

	return &job
}

func RemoveJob(id int) {
	jobs.Mutex.Lock()
	defer jobs.Mutex.Unlock()

	delete(jobs.jobMap, id)
}

func RemoveResult(id int) {
	jobs.Mutex.Lock()
	defer jobs.Mutex.Unlock()

	delete(jobs.resultMap, id)
}

func CreateJob(acquireFunc func() (string, error)) *Job {
	job := NewJob()
	go func() {
		s, err := acquireFunc()
		result := Result{
			Output: s,
			Err:    err,
		}
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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		web.JSONResponseError(errors.New("invalid id"), w)
	}
	if job, ok := jobs.jobMap[id]; ok {
		select {
		case result := <-job.ResultChannel:
			jobs.resultMap[id] = result
			RemoveJob(id)
			web.JSONResponse(
				web.StatusResponse{
					Status: "complete",
					Link:   "/api/v1/_jobs/result/" + strconv.Itoa(id),
				},
				w)
		default:
			web.JSONResponse(web.StatusResponse{Status: "inprogress"}, w)
		}
	} else {
		web.JSONResponseError(errors.New("not found"), w)
	}
}

func routerAcquireResult(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		web.JSONResponseError(errors.New("invalid id"), w)
	}
	if result, ok := jobs.resultMap[id]; ok {
		if result.Err != nil {
			web.JSONResponseError(result.Err, w)
		} else {
			if !validator.IsEmpty(result.Output) {
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
		web.JSONResponseError(errors.New("not found"), w)
	}
}

func RegisterRouterJobs(router *mux.Router) {
	jobs = New()

	n := router.PathPrefix("/_jobs").Subrouter().StrictSlash(false)

	n.HandleFunc("/status/{id}", routerAcquireStatus).Methods("GET")
	n.HandleFunc("/result/{id}", routerAcquireResult).Methods("GET")
}
