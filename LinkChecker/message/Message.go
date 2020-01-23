package message

import "encoding/json"

// Message is a simple struct containing the data of a incoming 'links to check' kafka message.
type Message struct {
	CommentTime int64
	CommentID string
	AuthorID string
	AuthorName string
	CommentText string
	Upvotes int
	PostID string
	PostTitle string
	SubredditID string
	SubredditTitle string
	ParentCommentID string
	IsReply bool
	NumReplies int
	ReplyIds []string
	LinksToCheck []string

	KafkaKey []byte `json:"-"`
}
// ToJSONWithoutLinks converts a message object back into json-encoded UTF-8 bytes, minus the 'links to check' field
func (m *Message) ToJSONWithoutLinks() ([]byte, error) {
	data := struct {
		CommentTime int64 `json:"commentTime"`
		CommentID string `json:"commentId"`
		AuthorID string `json:"authorId"`
		AuthorName string `json:"authorName"`
		CommentText string `json:"commentText"`
		Upvotes int `json:"upvotes"`
		PostID string `json:"postId"`
		PostTitle string `json:"postTitle"`
		SubredditID string `json:"subredditId"`
		SubredditTitle string `json:"subredditTitle"`
		ParentCommentID string `json:"parentCommentId"`
		IsReply bool `json:"isReply"`
		NumReplies int `json:"numReplies"`
		ReplyIds []string `json:"replyIds"`
	}{
		CommentTime: m.CommentTime,
		CommentID: m.CommentID,
		AuthorID: m.AuthorID,
		AuthorName: m.AuthorName,
		CommentText: m.CommentText,
		Upvotes: m.Upvotes,
		PostID: m.PostID,
		PostTitle: m.PostTitle,
		SubredditID: m.SubredditID,
		SubredditTitle: m.SubredditTitle,
		ParentCommentID: m.ParentCommentID,
		IsReply: m.IsReply,
		NumReplies: m.NumReplies,
		ReplyIds: m.ReplyIds,
	};

	return json.Marshal(&data);
}
// BuildMessage converts an incoming json string into a 'Message' object for use with this code-package
func BuildMessage(jsonWithLinks []byte, key []byte) (*Message, error) {
	var m Message;
	if err := json.Unmarshal(jsonWithLinks, &m); err != nil {
		return nil, err;
	}
	m.KafkaKey = key;
	return &m, nil;
}