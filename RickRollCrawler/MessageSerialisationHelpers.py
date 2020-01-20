# Functions which define the JSON formatting for each type of 'message' which will enqueued into our Kafka topics.
# This must match the schema which our consumers expect. We will return a JSON-string which is subsequently serialised into
# bytes for the sake of kafka
import json
from CommentAnalysisHelpers import getRedirectLinks

def newRickRollSerialiser(comment):
    return json.dumps(getObj(comment))

def redirectLinkSerialiser(comment):
    d = getObj(comment)
    d['linksToCheck'] = getRedirectLinks(comment.body)
    return json.dumps(d)

# Private
def getObj(comment):
    isReply = comment.link_id in comment.parent_id
    parId = comment.parent_id if isReply else ""

    obj = {
        "commentTime": comment.created_utc,
        "commentId": comment.id,
        "authorId": comment.author.id,
        "authorName": comment.author.name,
        "commentText": comment.body,
        "upvotes": comment.score,
        "postId": comment.link_id,
        "postTitle": comment.submission.title,
        "subredditId": comment.subreddit_id,
        "subredditTitle": comment.subreddit.display_name,
        "parentCommentId": parId,
        "isReply": isReply,
        "numReplies": len(comment.replies), # For now, we only track top-level replies
        "replyIds": list(map(lambda r: r.id, comment.replies))
    }
    return obj