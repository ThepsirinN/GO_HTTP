package handlersV1

import (
	"go_http_barko/utility/logger"
	"net/http"
)

func (h *handlersV1) GetHelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "EiEi")
	switch r.Method {
	case http.MethodGet:
		serviceData := h.svc.GetHelloWorldSvc(ctx)
		w.Write([]byte(serviceData))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
