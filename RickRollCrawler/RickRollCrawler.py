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
from CommentAnalysisHelpers import isRickRoll, getLinks
from MessageSerialisationHelpers import newRickRollSerialiser, redirectLinkSerialiser
import datetime
import os

# Reddit API Keys and details, loaded from the environment
CLIENTID = os.environ.get("PRAW_CLIENT_ID")
CLIENTSECRET = os.environ.get("PRAW_CLIENT_SECRET")
USERAGENT = "RickRollRankings:v0.1 (by u/AlexKfridges)"

# Kafka bootstrap server address and topic-names to enqueue, loaded from the environment
KAFKA_BOOTSTRAP_SERV = os.environ.get("KAFKA_BOOTSTRAP_SERV")
RICKROLL_TOPIC = os.environ.get("KAFKA_NEW_RICKROLL_TOPIC_NAME")
REDIRECT_TOPIC = os.environ.get("KAFKA_CHECK_REDIRECT_TOPIC_NAME")

# Instantiate a Kafka producer object, which connects to our cluster using the bootstrap server address.
producer = KafkaProducer(
    client_id = "PrawCrawler",
    bootstrap_servers = [KAFKA_BOOTSTRAP_SERV], 
    key_serializer = lambda s : s.encode(encoding='UTF-8',errors='strict'),
    value_serializer = lambda s : s.encode(encoding='UTF-8',errors='strict'),
    acks = 1
)

# Instantiate a 'Reddit' object, which acts as a logged-in conduit to reddit :)
redditApi = praw.Reddit(client_id = CLIENTID, client_secret = CLIENTSECRET, user_agent = USERAGENT)

print("Logged in as:")
print(redditApi.user.me())

total = 0

try:
    while True:
        try:
            for comment in redditApi.subreddit('all').stream.comments():
                total = total + 1
                if total % 1000 == 0:
                    print("Parsed " + str(total) + " comments")
                
                body = comment.body
                if (isRickRoll(body)):
                    # Publish a new 'New Rick-roll comment' event!
                    print("Found rick roll")
                    producer.send(
                        topic = RICKROLL_TOPIC,
                        value = newRickRollSerialiser(comment),
                        key = comment.author.id
                    )
                else:
                    links = getLinks(body)
                    if (len(links) > 0):
                        # Publish a new 'redirect links to check' event!
                        print("Found redirect link candidate:")
                        print(links)
                        producer.send(
                            topic = REDIRECT_TOPIC,
                            value = redirectLinkSerialiser(comment, links),
                            key = comment.author.id
                        )
        except praw.exceptions.APIException:
            print(" error!\nRate limited on comment", comment.id, "by", comment.author)
        except prawcore.exceptions.Forbidden:
            print(" error!\n403 error on comment", comment.id, "by", comment.author)
        except prawcore.exceptions.ServerError:
            print(" error!\nServer error on comment")
        except UnicodeDecodeError:
            print(" error!\nCould not encode comment using UTF-8 string encoding")
        except Exception:
            print(" CAUGHT GENERIC ERROR!! THIS IS BAD, BUT JUST DOING IT FOR TESTING")
finally:
    producer.flush(5)
    producer.close()