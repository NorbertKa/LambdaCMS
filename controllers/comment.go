package controller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/NorbertKa/LambdaCMS/jwt"
	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/julienschmidt/httprouter"
)

func (h Handler) Comment_POST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
	commentId := r.Form.Get("commentId")
	commentId_int, err := strconv.Atoi(commentId)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	var comId sql.NullInt64
	rawComment := db.Comment{
		UserId:    decodedToken.UserId,
		PostId:    postId_int,
		Body:      body,
		CommentId: comId,
	}
	if commentId_int > 0 {
		checkComment, _ := db.Comment_GetById(h.DB, commentId_int)
		if checkComment.Id <= 0 {
			response := Response{
				Status:  false,
				Message: "Nonexistant Top comment",
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
		} else {
			err := rawComment.CommentId.Scan(checkComment.Id)
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(ErrResponseInternalServerError))
				return
			}
		}
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

func (h Handler) Comment_UPVOTE(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response    Response
		CommentVote *db.CommentVote `json:"commentVote,omitempty"`
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
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		return
	}
	commentId := ps.ByName("ID")
	commentId_int, err := strconv.Atoi(commentId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	comment, err := db.Comment_GetById(h.DB, commentId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comment.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Comment not found",
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
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		return
	}
	checkCommentVote, err := db.CommentVote_GetByUserData(h.DB, decodedToken.UserId, commentId_int)
	if err != nil {
		fmt.Println(err)
	}
	if checkCommentVote.Type == "upvote" {
		response := Response{
			Status:  false,
			Message: "Already upvoted",
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
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		return
	} else if checkCommentVote.Type == "downvote" {
		err := checkCommentVote.Change(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		err = comment.RemoveDownvote(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		err = comment.Upvote(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		response := Response{
			Status:  false,
			Message: "Upvoted from Downvote",
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
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		return
	} else {
		rawCommentVote := db.CommentVote{
			UserId:    decodedToken.UserId,
			CommentId: commentId_int,
			Type:      "upvote",
		}
		err := rawCommentVote.Create(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		err = comment.Upvote(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		response := Response{
			Status:  false,
			Message: "Upvoted",
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
}

func (h Handler) Comment_DOWNVOTE(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response    Response
		CommentVote *db.CommentVote `json:"commentVote,omitempty"`
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
	commentId := ps.ByName("ID")
	commentId_int, err := strconv.Atoi(commentId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	comment, err := db.Comment_GetById(h.DB, commentId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comment.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Comment not found",
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
	checkCommentVote, _ := db.CommentVote_GetByUserData(h.DB, decodedToken.UserId, commentId_int)
	if checkCommentVote.Type == "downvote" {
		response := Response{
			Status:  false,
			Message: "Already downvoted",
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
	} else if checkCommentVote.Type == "upvote" {
		err := checkCommentVote.Change(h.DB)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		err = comment.RemoveUpvote(h.DB)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		err = comment.Downvote(h.DB)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		response := Response{
			Status:  false,
			Message: "Downvoted from Upvote",
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
		rawCommentVote := db.CommentVote{
			UserId:    decodedToken.UserId,
			CommentId: commentId_int,
			Type:      "downvote",
		}
		err := rawCommentVote.Create(h.DB)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseUnsupportedContentType))
			return
		}
		response := Response{
			Status:  false,
			Message: "Downvoted",
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
}

func (h Handler) Comment_GET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Comment  *db.Comment `json:"comment,omitempty"`
	}
	commentId := ps.ByName("ID")
	commentId_int, err := strconv.Atoi(commentId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	comment, err := db.Comment_GetById(h.DB, commentId_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comment.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Comment not found",
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
		Message: "Comment found",
	}
	resp := Resp{
		Response: response,
		Comment:  &comment,
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

func (h Handler) Comment_DELETE(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Comment  *db.Comment `json:"comment,omitempty"`
	}
	commentId := ps.ByName("ID")
	commentId_int, err := strconv.Atoi(commentId)
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
	comment, err := db.Comment_GetById(h.DB, commentId_int)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comment.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Comment not found",
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
	if comment.UserId != decodedToken.UserId {
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
		err = comment.Delete(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		response := Response{
			Status:  true,
			Message: "Comment Deleted",
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

func (h Handler) Comment_UPDATE(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	type Resp struct {
		Response Response
		Comment  *db.Comment `json:"comment,omitempty"`
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
	commentID := ps.ByName("ID")
	commentID_int, err := strconv.Atoi(commentID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	comment, err := db.Comment_GetById(h.DB, commentID_int)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	if comment.Id == 0 {
		response := Response{
			Status:  false,
			Message: "Comment not found",
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
			Message: "Invalid Comment body",
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
	if comment.UserId != decodedToken.UserId {
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
		comment.Body = newBody
		err = comment.Update(h.DB)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(ErrResponseInternalServerError))
			return
		}
		response := Response{
			Status:  true,
			Message: "Comment updated",
		}
		resp := Resp{
			Response: response,
			Comment:  &comment,
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

func (h Handler) Comment_GETComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type Resp struct {
		Response Response
		Comments db.Comments
	}
	commentID := ps.ByName("ID")
	commentID_int, err := strconv.Atoi(commentID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(ErrResponseInternalServerError))
		return
	}
	comments, err := db.Comment_GetByCommentId(h.DB, commentID_int)
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
