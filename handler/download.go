package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"the-rush/graph/model"
	"the-rush/repository"
)

type DownloadHandler struct {
	repo *repository.Repository
}

func NewDownloadHandler(repo *repository.Repository) *DownloadHandler {
	return &DownloadHandler{repo: repo}
}

func (x *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("downloadCSVHandler.ServeHTTP")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("download requires POST with body of players filtering/pagination arguments"))
		return
	}

	var (
		args = new(model.PlayersArgs)
		ctx  = r.Context()
	)
	err = json.Unmarshal(body, args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	page, err := model.PlayerPaginationFromArgs(*args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	resp, err := x.repo.Player().Search(ctx, args.Name, *page, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	csv := model.PlayersCSV(resp.Players).Marshal()
	w.Header().Set("Content-Disposition", "attachment; filename=players.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(csv)))
	_, _ = w.Write(csv)
}
