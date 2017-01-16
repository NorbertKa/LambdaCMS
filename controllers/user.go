package controller

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/julienschmidt/httprouter"
)

func (h Handler) User_Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("ID")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	user, err := db.User_GetById(h.DB, id_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		User     db.User
	}
	response := Response{
		Status:  true,
		Message: "User found",
	}
	user.Password = ""
	resp := Resp{
		Response: response,
		User:     user,
	}
	ContentType := r.Header.Get("Response-Content-Type")
	if ContentType == "" || ContentType == "application/json" {
		js, err := json.Marshal(resp)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if ContentType == "application/xml" {
		x, err := xml.MarshalIndent(resp, "", "  ")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(x)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseUnsupportedContentType))
		return
	}
}

func (h Handler) User_Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response `json:"response"`
		User     db.User  `json:"user,omitempty"`
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
	user := db.User{
		Username: username,
		Password: password,
		Role:     "user",
	}
	err = user.Validate()
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
			ErrCode: http.StatusBadRequest,
		}
		resp := Resp{
			Response: response,
			User:     db.User{},
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	} else {
		err := user.Create(h.DB)
		if err != nil {
			response := Response{
				Status:  false,
				Message: "User already exists",
				ErrCode: http.StatusBadRequest,
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
		response := Response{
			Status:  true,
			Message: "User created",
		}
		realUser, err := db.User_GetByUsername(h.DB, user.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		realUser.Hash = ""
		realUser.Password = ""
		resp := Resp{
			Response: response,
			User:     realUser,
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	}
}

func (h Handler) Users_Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response
		Users    db.Users
	}
	users, err := db.User_GetAll(h.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if users == nil {
		response := Response{
			Status:  false,
			Message: "No users found",
		}
		resp := Resp{
			Response: response,
			Users:    db.Users{},
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	} else {
		response := Response{
			Status:  true,
			Message: "List of Users",
		}
		resp := Resp{
			Response: response,
			Users:    users,
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
	}
}
