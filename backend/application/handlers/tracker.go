package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"pento/code-challenge/domain/tracker/models"
	"pento/code-challenge/domain/tracker/services"
	"pento/code-challenge/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TrackerService interface {
	GetTracker(ctx context.Context, id uint64) (models.TimeTracker, error)
	ListTrackers(ctx context.Context, params services.ListTimeTracker) ([]models.TimeTracker, error)
	CreateTracker(ctx context.Context, params services.CreateTrackerParams) (models.TimeTracker, error)
	UpdateTracker(ctx context.Context, params services.UpdateTrackerParams) (models.TimeTracker, error)
	DeleteTracker(ctx context.Context, params services.DeleteTrackerParams) error
}

type TrackerHandler struct {
	service TrackerService
}

func NewTrackerHandler(service TrackerService) *TrackerHandler {
	return &TrackerHandler{
		service: service,
	}
}

type updateTimeTrackerRequest struct {
	End     time.Time `json:"end"`
	Name    string    `json:"name"`
	Version uint32    `json:"version"`
}

type createTrackerRequest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Name  string    `json:"name"`
}

type TimeTrackerResponse struct {
	ID        *uint64    `json:"id"`
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
	Name      *string    `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   uint32     `json:"version"`
}

type TimeTrackersResponse struct {
	Trackers []TimeTrackerResponse `json:"trackers"`
}

func (h TrackerHandler) GetTracker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	Tracker, err := h.service.GetTracker(context.Background(), i)
	if err != nil {
		switch err {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}

	response, err := json.Marshal(fromDomain(Tracker))
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

func (h TrackerHandler) ListTrackers(w http.ResponseWriter, r *http.Request) {

	var startDate, endDate time.Time
	var err error

	if r.FormValue("start_date") == "" && r.FormValue("end_date") == "" {
		startDate = time.Time{}
		endDate = time.Time{}
	} else {
		startDate, err = utils.StrToTime(r.FormValue("start_date"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}
		endDate, err = utils.StrToTime(r.FormValue("end_date"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}
	}

	trackers, err := h.service.ListTrackers(context.Background(), services.ListTimeTracker{
		Start: startDate,
		End:   endDate,
	})
	if err != nil {
		switch err {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}

	response, err := json.Marshal(fromDomainSlice(trackers))
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

func (h TrackerHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {

	var request createTrackerRequest
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
		Name:  request.Name,
	}

	Tracker, err := h.service.CreateTracker(context.Background(), params)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(fromDomain(Tracker))
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

func (h TrackerHandler) UpdateTracker(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	paramID := vars["id"]

	id, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
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
		Name:    request.Name,
		End:     request.End,
		ID:      id,
	}

	tracker, err := h.service.UpdateTracker(context.Background(), params)
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

	response, err := json.Marshal(fromDomain(tracker))
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

func (h TrackerHandler) DeleteTracker(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	paramsID := vars["id"]

	id, err := strconv.ParseUint(paramsID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)

		return
	}

	err = h.service.DeleteTracker(context.Background(), services.DeleteTrackerParams{
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

	var end *time.Time = nil

	if !tracker.End.IsZero() {
		end = &tracker.End
	}

	return TimeTrackerResponse{
		ID:        &tracker.ID,
		Start:     &tracker.Start,
		End:       end,
		Name:      &tracker.Name,
		CreatedAt: tracker.Meta.GetCreatedAt(),
		UpdatedAt: tracker.Meta.GetUpdatedAt(),
		Version:   tracker.Meta.GetVersion(),
	}
}

func fromDomainSlice(Trackers []models.TimeTracker) []TimeTrackerResponse {
	var timeTrackerResponse []TimeTrackerResponse = make([]TimeTrackerResponse, 0)

	for _, elem := range Trackers {
		u := fromDomain(elem)

		timeTrackerResponse = append(timeTrackerResponse, u)
	}

	return timeTrackerResponse
}
