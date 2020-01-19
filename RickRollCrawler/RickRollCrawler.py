# This is a simple reddit crawling bot which subscribes to Reddit's /r/all comment feed API endpoint.
# We will simply query for ALL possible comments every couple of seconds, and scan through the comments for
# rick-roll links. If we find any, we will publish the comment as a 'new rick roll' event into Kafka for
# processing by our system.

# praw is the python module that makes it easy to interface with the reddit api.
import praw
import prawcore

# kafka-python is a simple library that allows us to behave as a kafka producer.


import datetime
import os



# Test app client id, client secret, and pass word, etc.
CLIENTID = os.environ.get("CLIENT_ID")
CLIENTSECRET = os.environ.get("CLIENT_SECRET")
USERAGENT = "RickRollRankings:v0.1 (by u/AlexKfridges)"

# List of url-snippets that are known rick-roll videos
rickRollTextSnippets = [
    # Youtube videos
    "oHg5SJYRHA0",
    "dQw4w9WgXcQ",
    "xfr64zoBTAQ",
    "dPmZqsQNzGA",
    "r8tXjJL3xcM",
    "6-HUgzYPm9g",

    # Fake news websites which contains an embedded Rick-roll on the site
    "latlmes.com/breaking",
    "rickrolled.com"
]

# List of common link-shortening services which should be checked for rick-roll-redirects
redirectCandidateSnippets = [
    "tinyurl.com/",
    "bit.ly/",
    "bitly.com/",
    "alturl.com/",
    "goo.gl/",
    "rb.gy/",
    "bit.do/",
    "is.gd/",
    "ow.ly/"
]

# Instantiate a 'Reddit' object, which acts as a logged-in conduit to reddit :)
redditApi = praw.Reddit(client_id = CLIENTID, client_secret = CLIENTSECRET, user_agent = USERAGENT)

print("Logged in as:")
print(redditApi.user.me())

def isRickRoll(text):
    if "oHg5SJYRHA0".lower() in text.lower() or "dQw4w9WgXcQ".lower() in text.lower():
        return True
    else:
        return False

counter = 0
totalCounter = 0

# Open a file in the current directory for saving rick-roll references to.
with open('d:/RickRollRankingsV1/RickRollCrawler/datalogs/commentstracking.txt', 'a+') as datafile:
    while True:
        try:
            for comment in redditApi.subreddit('all').stream.comments():
                totalCounter = totalCounter + 1
                if totalCounter % 5000 == 0:
                    print(" - - - - > PARSED " + str(totalCounter) + " COMMENTS < - - - - ")
                if (isRickRoll(comment.body)):
                    counter = counter + 1
                    datafile.write("Rickroll #" + str(counter) + " (" + str(totalCounter) + ") | " + datetime.datetime.fromtimestamp(comment.created_utc).strftime('%Y-%m-%d %H:%M:%S') + " | Comment id: " + comment.id + " by " + comment.author.name + " | Thread: " + comment.subreddit.display_name + " - " + comment.submission.title + "\n")
                    print("================ RICKROLL #" + str(counter) + " ================")
                    print("Thread: " + comment.subreddit.display_name + " -> " + comment.submission.title)
                    print("-----------------------------------------------")
                    print("Comment: (" + str(comment.score) + ") " + comment.author.name + " says: ")
                    print(comment.body)
                    print()
        except praw.exceptions.APIException:
            print(" error!\nRate limited on comment", comment.id, "by", comment.author)
        except prawcore.exceptions.Forbidden:
            print(" error!\n403 error on comment", comment.id, "by", comment.author)
        except prawcore.exceptions.ServerError:
            print(" error!\nServer error on comment")
