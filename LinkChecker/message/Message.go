package message

import "encoding/json"

// Message is a simple struct containing the data of a incoming 'links to check' kafka message.
type Message struct {
	CommentTime int
	CommentId string
	AuthorId string
	AuthorName string
	CommentText string
	Upvotes int
	PostId string
	PostTitle string
	SubredditId string
	SubredditTitle string
	ParentCommentId string
	IsReply bool
	NumReplies int
	ReplyIds []string
	LinksToCheck []string
}
// ToJSONWithoutLinks converts a message object back into json-encoded UTF-8 bytes, minus the 'links to check' field
func (m *Message) ToJSONWithoutLinks() ([]byte, error) {
	data := struct {
		commentTime int
		commentId string
		authorId string
		authorName string
		commentText string
		upvotes int
		postId string
		postTitle string
		subredditId string
		subredditTitle string
		parentCommentId string
		isReply bool
		numReplies int
		replyIds []string
	}{
		commentTime: m.CommentTime,
		commentId: m.CommentId,
		authorId: m.AuthorId,
		authorName: m.AuthorName,
		commentText: m.CommentText,
		upvotes: m.Upvotes,
		postId: m.PostId,
		postTitle: m.PostTitle,
		subredditId: m.SubredditId,
		subredditTitle: m.SubredditTitle,
		parentCommentId: m.ParentCommentId,
		isReply: m.IsReply,
		numReplies: m.NumReplies,
		replyIds: m.ReplyIds,
	};
	return json.Marshal(data);
}
// BuildMessage converts an incoming json string into a 'Message' object for use with this code-package
func BuildMessage(jsonWithLinks []byte) (*Message, error) {
	var m Message;
	if err := json.Unmarshal(jsonWithLinks, &m); err != nil {
		return nil, err;
	}
	return &m, nil;
}