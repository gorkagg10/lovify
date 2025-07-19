package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
	"net/http"
	"time"
)

type SendMessageRequest struct {
	Content string `json:"content"`
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	matchID := params["match_id"]
	var sendMessageRequest SendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&sendMessageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.MessagingServiceClient.SendMessage(
		r.Context(),
		&messagingServiceGrpc.SendMessageRequest{
			MatchID: &matchID,
			UserID:  &userID,
			Content: &sendMessageRequest.Content,
		})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type ListMessagesResponse struct {
	MatchID    string `json:"match_id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Content    string `json:"content"`
	SendAt     string `json:"send_at"`
	Read       bool   `json:"read"`
}

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	matchID := params["match_id"]

	messageList, err := h.MessagingServiceClient.ListMessages(
		r.Context(),
		&messagingServiceGrpc.ListMessagesRequest{
			UserID:  &userID,
			MatchID: &matchID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	listMessagesResponse := make([]ListMessagesResponse, len(messageList.Messages))
	for i, message := range messageList.Messages {
		listMessagesResponse[i] = ListMessagesResponse{
			MatchID:    matchID,
			FromUserID: message.GetFromUserID(),
			ToUserID:   message.GetToUserID(),
			Content:    message.GetContent(),
			SendAt:     message.GetSendAt().AsTime().Format(time.RFC3339),
			Read:       message.GetRead(),
		}
	}

	listMessagesResponseJSON, err := json.Marshal(listMessagesResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(listMessagesResponseJSON)
}
