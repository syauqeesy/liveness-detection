package handler

import (
	"net/http"
)

type inferenceHandler handler

func (h *inferenceHandler) Execute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	file, header, err := r.FormFile("image")
	if err != nil {
		h.CommonHttp.ErrorHandler(w, err, nil)
		return
	}

	result, err := h.Service.Inference.Execute(r.Context(), file, header)
	if err != nil {
		h.CommonHttp.ErrorHandler(w, err, nil)
		return
	}

	h.CommonHttp.WriteResponse(w, http.StatusOK, "Liveness checked!", result)
}
