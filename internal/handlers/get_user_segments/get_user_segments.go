package get_user_segments

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type segmentsGetter interface {
	GetUserSegments(context.Context, uint64) ([]string, error)
}

func New(service segmentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["userID"]
		userID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid user id"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		segments, err := service.GetUserSegments(ctx, userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Content-Type", "encoding/json")
		if err := json.NewEncoder(w).Encode(segments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
