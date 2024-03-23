package handlersV1

import (
	"net/http"
)

func (h *handlersV1) GetHelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		serviceData := h.svc.GetHelloWorldSvc(h.ctx)
		w.Write([]byte(serviceData))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
