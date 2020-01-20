# Some helper functions to check if a comment contains a rick-roll link, or a known redirect link
from urlextract import URLExtract

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

def isRickRoll(text):
    return any(linktext in text for linktext in rickRollTextSnippets)

def isRedirectLink(text):
    return any(linktext in text for linktext in redirectCandidateSnippets)

def getRedirectLinks(text):
    extractor = URLExtract()
    possibleUrls = extractor.find_urls(text)
    matchesList = list(filter(lambda s: any(snippet in s for snippet in redirectCandidateSnippets), possibleUrls))
    print("Found possible redirect links in comment: ")
    print(matchesList)
    return matchesList