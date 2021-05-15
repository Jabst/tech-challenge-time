package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"pento/code-challenge/domain/tracker/models"
	"pento/code-challenge/domain/tracker/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TrackerService interface {
	GetUser(ctx context.Context, id int) (models.TimeTracker, error)
	ListUsers(ctx context.Context, queryTerms map[string]string) ([]models.TimeTracker, error)
	CreateUser(ctx context.Context, params services.CreateTrackerParams) (models.TimeTracker, error)
	UpdateUser(ctx context.Context, params services.UpdateTrackerParams) (models.TimeTracker, error)
	DeleteUser(ctx context.Context, params services.DeleteTrackerParams) error
}

type UserHandler struct {
	service TrackerService
}

func NewUserHandler(service TrackerService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type updateTimeTrackerRequest struct {
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Name    string    `json:"name"`
	ID      uint64    `json:"id"`
	Version uint32    `json:"version"`
}

type createUserRequest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Name  string    `json:"name"`
}

type TimeTrackerResponse struct {
	ID        uint64    `json:"id"`
	Start     string    `json:"start"`
	End       string    `json:"end"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   uint32    `json:"version"`
}

type TimeTrackersResponse struct {
	Trackers []TimeTrackerResponse `json:"trackers"`
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}

	user, err := h.service.GetUser(context.Background(), i)
	if err != nil {
		switch err {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}

	response, err := json.Marshal(fromDomain(user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	_, err = w.Write(response)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var request createUserRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	params := services.CreateTrackerParams{
		Start: request.Start,
		End:   request.End,
		Name:  request.Name,
	}

	user, err := h.service.CreateUser(context.Background(), params)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(fromDomain(user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}
}

func (h UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	paramID := vars["id"]

	id, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}

	var request updateTimeTrackerRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	params := services.UpdateTrackerParams{
		Version: request.Version,

		ID: id,
	}

	user, err := h.service.UpdateUser(context.Background(), params)
	if err != nil {
		switch err {
		case services.ErrTrackerNotFound:
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		case services.ErrWrongVersion:
			w.WriteHeader(http.StatusConflict)
			log.Println(err)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	response, err := json.Marshal(fromDomain(user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func (h UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	paramsID := vars["id"]

	id, err := strconv.ParseUint(paramsID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	err = h.service.DeleteUser(context.Background(), services.DeleteTrackerParams{
		ID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func fromDomain(tracker models.TimeTracker) TimeTrackerResponse {
	return TimeTrackerResponse{
		ID:        tracker.ID,
		Start:     tracker.Start.String(),
		End:       tracker.End.String(),
		Name:      tracker.Name,
		CreatedAt: tracker.Meta.GetCreatedAt(),
		UpdatedAt: tracker.Meta.GetUpdatedAt(),
		Version:   tracker.Meta.GetVersion(),
	}
}

func fromDomainSlice(users []models.TimeTracker) TimeTrackersResponse {
	var timeTrackerResponse []TimeTrackerResponse = make([]TimeTrackerResponse, 0)

	for _, elem := range users {
		u := fromDomain(elem)

		timeTrackerResponse = append(timeTrackerResponse, u)
	}

	return TimeTrackersResponse{
		Trackers: timeTrackerResponse,
	}
}
