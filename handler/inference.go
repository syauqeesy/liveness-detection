package handler

import (
	"net/http"

	"github.com/syauqeesy/liveness-detection/common"
)

type inferenceHandler handler

func (h *inferenceHandler) Predict(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := common.ContextWithRequestId(r.Context(), w.Header().Get("X-Request-Id"))

	file, header, err := r.FormFile("image")
	if err != nil {
		h.CommonHttp.ErrorHandler(w, err, nil)
		return
	}

	mode := r.FormValue("mode")
	if mode != "managed_service" && mode != "self_managed_service" {
		h.CommonHttp.ErrorHandler(w, err, nil)
		return
	}

	result, err := h.Service.Inference.Predict(ctx, mode, file, header)
	if err != nil {
		h.CommonHttp.ErrorHandler(w, err, nil)
		return
	}

	h.CommonHttp.WriteResponse(w, http.StatusOK, "Liveness checked!", result)
}
