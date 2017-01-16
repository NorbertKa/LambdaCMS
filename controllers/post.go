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

func (h Handler) Post_POST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		Post     *db.Post `json:"board,omitempty"`
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
	body := r.Form.Get("body")
	boardId := r.Form.Get("boardId")
	boardId_int, err := strconv.Atoi(boardId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseUnsupportedContentType))
		return
	}
	checkBoard, _ := db.Board_GetById(h.DB, boardId_int)
	if checkBoard.Id <= 0 {
		response := Response{
			Status:  false,
			Message: "Unexistant board",
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
	rawPost := db.Post{
		UserId:  decodedToken.UserId,
		BoardId: boardId_int,
		Title:   title,
		Body:    body,
	}
	err = rawPost.Create(h.DB)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	response := Response{
		Status:  true,
		Message: "Post created",
	}
	resp := Resp{
		Response: response,
		Post:     &rawPost,
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

func (h Handler) Post_GET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Post     *db.Post `json:"post,omitempty"`
	}
	postId := ps.ByName("ID")
	postId_int, err := strconv.Atoi(postId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	post, err := db.Post_GetById(h.DB, postId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if post.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Post not found",
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
		Message: "Post found",
	}
	resp := Resp{
		Response: response,
		Post:     &post,
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

func (h Handler) Post_Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Post     *db.Post `json:"post,omitempty"`
	}
	postId := ps.ByName("ID")
	postId_int, err := strconv.Atoi(postId)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
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
	post, err := db.Post_GetById(h.DB, postId_int)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if post.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Post not found",
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
	if post.UserId != decodedToken.UserId {
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
	} else {
		err = post.Delete(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		response := Response{
			Status:  true,
			Message: "Post Deleted",
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
	}
}

func (h Handler) Post_Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		Post     *db.Post `json:"post,omitempty"`
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
	postId := ps.ByName("ID")
	postId_int, err := strconv.Atoi(postId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	post, err := db.Post_GetById(h.DB, postId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if post.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Post not found",
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
	newBody := r.Form.Get("body")
	if len(newBody) <= 0 {
		response := Response{
			Status:  false,
			Message: "Invalid Post body",
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
	if post.UserId != decodedToken.UserId {
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
	} else {
		post.Body = newBody
		err = post.Update(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		response := Response{
			Status:  true,
			Message: "Post found",
		}
		resp := Resp{
			Response: response,
			Post:     &post,
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
}

func (h Handler) Posts_GET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Resp struct {
		Response Response
		Posts    db.Posts
	}
	posts, err := db.Post_GetAll(h.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if posts == nil {
		response := Response{
			Status:  false,
			Message: "No posts found",
		}
		resp := Resp{
			Response: response,
			Posts:    db.Posts{},
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
			Message: "List of Posts",
		}
		resp := Resp{
			Response: response,
			Posts:    posts,
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
