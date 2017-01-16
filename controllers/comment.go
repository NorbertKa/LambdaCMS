package controller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NorbertKa/LambdaCMS/jwt"
	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/julienschmidt/httprouter"
)

func (h Handler) Comment_POST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		Comment  *db.Comment `json:"board,omitempty"`
	}
	token := r.Header.Get("token")
	decodedToken, err := jwt.DecodeToken(token, h.Conf.Secret)
	if err != nil {
		response := Response{
			Status:  false,
			Message: "Invalid token",
		}
		resp := Resp{
			Response: response,
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Println(err)
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
		return
	}
	body := r.Form.Get("body")
	postId := r.Form.Get("postId")
	postId_int, err := strconv.Atoi(postId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseUnsupportedContentType))
		return
	}
	checkPost, _ := db.Post_GetById(h.DB, postId_int)
	if checkPost.Id <= 0 {
		response := Response{
			Status:  false,
			Message: "Unexistant post",
		}
		resp := Resp{
			Response: response,
		}
		ContentType := r.Header.Get("Response-Content-Type")
		if ContentType == "" || ContentType == "application/json" {
			js, err := json.Marshal(resp)
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else if ContentType == "application/xml" {
			x, err := xml.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Println(err)
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
	rawComment := db.Comment{
		UserId: decodedToken.UserId,
		PostId: postId_int,
		Body:   body,
	}
	err = rawComment.Create(h.DB)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	response := Response{
		Status:  true,
		Message: "Comment created",
	}
	resp := Resp{
		Response: response,
		Comment:  &rawComment,
	}
	ContentType := r.Header.Get("Response-Content-Type")
	if ContentType == "" || ContentType == "application/json" {
		js, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	} else if ContentType == "application/xml" {
		x, err := xml.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Println(err)
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

func (h Handler) Comments_GET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response
		Comments db.Comments
	}
	comments, err := db.Comment_GetAll(h.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comments == nil {
		response := Response{
			Status:  false,
			Message: "No comments found",
		}
		resp := Resp{
			Response: response,
			Comments: db.Comments{},
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
			Message: "List of Coments",
		}
		resp := Resp{
			Response: response,
			Comments: comments,
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
