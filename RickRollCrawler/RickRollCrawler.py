# This is a simple reddit crawling bot which subscribes to Reddit's /r/all comment feed API endpoint.
# We will simply query for ALL possible comments every couple of seconds, and scan through the comments for
# rick-roll links. If we find any, we will publish the comment as a 'new rick roll' event into Kafka for
# processing by our system.

# praw is the python module that makes it easy to interface with the reddit api.
import praw
import prawcore

# kafka-python is a simple library that allows us to behave as a kafka producer.
from kafka import KafkaProducer
from kafka.errors import KafkaError

# Misc other imports
from Helpers import isRickRoll, isRedirectLink
import datetime
import os

# Reddit API Keys and details, loaded from the environment
CLIENTID = os.environ.get("CLIENT_ID")
CLIENTSECRET = os.environ.get("CLIENT_SECRET")
USERAGENT = "RickRollRankings:v0.1 (by u/AlexKfridges)"

# Kafka bootstrap server address and topic-names to enqueue, loaded from the environment
KAFKA_BOOTSTRAP_SERV = os.environ.get("KAFKA_BOOTSTRAP_SERV")
RICKROLL_TOPIC = os.environ.get("KAFKA_NEW_RICKROLL_TOPIC_NAME")
REDIRECT_TOPIC = os.environ.get("KAFKA_CHECK_REDIRECT_TOPIC_NAME")
NUM_PRODUCER_RETIRES = os.environ.get("KAFKA_PRODUCER_NUM_RETRIES")

# Instantiate a Kafka producer object, which connects to our cluster using the bootstrap server address.
producer = KafkaProducer(bootstrap_servers=[KAFKA_BOOTSTRAP_SERV], retries=NUM_PRODUCER_RETIRES)

# Instantiate a 'Reddit' object, which acts as a logged-in conduit to reddit :)
redditApi = praw.Reddit(client_id = CLIENTID, client_secret = CLIENTSECRET, user_agent = USERAGENT)

print("Logged in as:")
print(redditApi.user.me())

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
