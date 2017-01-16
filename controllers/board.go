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

func (h Handler) Board_POST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		Board    *db.Board `json:"board,omitempty"`
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
	title := r.Form.Get("title")
	checkBoard, _ := db.Board_GetByTitle(h.DB, title)
	if checkBoard.Id > 0 {
		response := Response{
			Status:  false,
			Message: "Board already exists",
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
	miniDesc := r.Form.Get("miniDesc")
	fullDesc := r.Form.Get("fullDesc")
	image := r.Form.Get("image")
	rawBoard := db.Board{
		UserId:          decodedToken.UserId,
		Title:           title,
		MiniDescription: miniDesc,
		FullDescription: fullDesc,
		Image:           image,
	}
	err = rawBoard.Validate()
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
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
	err = rawBoard.Create(h.DB)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	response := Response{
		Status:  true,
		Message: "Board created",
	}
	board, err := db.Board_GetByTitle(h.DB, rawBoard.Title)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	resp := Resp{
		Response: response,
		Board:    board,
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

func (h Handler) Board_GET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Board    *db.Board `json:"board,omitempty"`
	}
	boardId := ps.ByName("ID")
	boardId_int, err := strconv.Atoi(boardId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	board, err := db.Board_GetById(h.DB, boardId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if board.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Board not found",
		}
		resp := Resp{
			Response: response,
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
			return
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
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		return
	}
	response := Response{
		Status:  true,
		Message: "Board found",
	}
	resp := Resp{
		Response: response,
		Board:    board,
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
		return
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
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseUnsupportedContentType))
		return
	}
}

func (h Handler) Boards_GET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response
		Boards   db.Boards
	}
	boards, err := db.Board_GetAll(h.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if boards == nil {
		response := Response{
			Status:  false,
			Message: "No boards found",
		}
		resp := Resp{
			Response: response,
			Boards:   db.Boards{},
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
			Message: "List of Boards",
		}
		resp := Resp{
			Response: response,
			Boards:   boards,
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