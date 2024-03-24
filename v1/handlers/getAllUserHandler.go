package handlersV1

import (
	"encoding/json"
	"go_http_barko/utility/logger"
	"net/http"
)

func (h *handlersV1) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		serviceData, err := h.svc.GetAllUserSvc(ctx)
		if err != nil {
			logger.Error(ctx, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			jsonData, _ := json.Marshal(serviceData)
			w.Write([]byte(jsonData))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
