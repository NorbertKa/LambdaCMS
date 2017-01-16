package controller

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/julienschmidt/httprouter"
)

func (h Handler) Login_POST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response `json:"response"`
	}
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	user, err := db.User_GetByUsername(h.DB, username)
	if err != nil {
		response := Response{
			Status:  false,
			Message: "Username or Password incorrect",
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(response)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(response, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	}
	check, _ := user.ComparePassword(h.DB, password)
	if !check {
		response := Response{
			Status:  false,
			Message: "Username or Password incorrect",
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(response)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(response, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	} else {
		response := Response{
			Status:  true,
			Message: "Logged in",
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(response)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(response, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	}
}
