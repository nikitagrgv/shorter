package http

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/nikitagrgv/shorter/internal/domain"
	"github.com/nikitagrgv/shorter/internal/usecase"
)

type LinkHandler struct {
	tmpl    *template.Template
	uc      *usecase.LinkUsecase
	baseUrl string
}

func NewLinkHandler(tmpl *template.Template, uc *usecase.LinkUsecase, baseUrl string) *LinkHandler {
	return &LinkHandler{tmpl: tmpl, uc: uc, baseUrl: baseUrl}
}

func (h *LinkHandler) ShowHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.tmpl.ExecuteTemplate(w, "hello", nil)
	if err != nil {
		h.handleError(w, err)
	}
}

func (h *LinkHandler) RedirectShortLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	link, err := h.uc.GetLinkByCode(r.Context(), code)
	if errors.Is(err, domain.ErrLinkNotFound) {
		h.handleError(w, err)
		return
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	orig := link.LongURL
	http.Redirect(w, r, orig, http.StatusFound)
}

func (h *LinkHandler) ShowResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	token := r.PathValue("token")
	link, err := h.uc.GetLinkByToken(r.Context(), token)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shortUrl, err := url.JoinPath(h.baseUrl, link.Short)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := ResultPageData{ShortURL: shortUrl}

	err = h.tmpl.ExecuteTemplate(w, "result", data)
	if err != nil {
		h.handleError(w, err)
	}
}

func (h *LinkHandler) PostLink(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.uc.CreateLink(r.Context(), usecase.CreateLinkCommand{LongURL: req.LongURL})
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := CreateLinkResponse{Code: res.AccessToken, ShortURL: res.Link.Short}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *LinkHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrLinkNotFound):
		h.render404(w)
	default:
		h.render500(w)
	}
}

func (h *LinkHandler) render404(w http.ResponseWriter) {
	data := ErrorPageData{ErrorCode: http.StatusNotFound, ErrorDescription: "Not Found"}
	h.renderError(w, data)
}

func (h *LinkHandler) render500(w http.ResponseWriter) {
	data := ErrorPageData{ErrorCode: http.StatusInternalServerError, ErrorDescription: "Internal Error"}
	h.renderError(w, data)
}

func (h *LinkHandler) renderError(w http.ResponseWriter, data ErrorPageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	err := h.tmpl.ExecuteTemplate(w, "error", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
