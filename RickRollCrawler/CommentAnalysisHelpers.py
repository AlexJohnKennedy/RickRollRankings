# Some helper functions to check if a comment contains a rick-roll link, or a known redirect link

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

def isRickRoll(comment):
    return any(linktext in comment.body for linktext in rickRollTextSnippets)

def isRedirectLink(comment):
    return any(linktext in comment.body for linktext in redirectCandidateSnippets)

def getRedirectLinks(comment):
    print("TODO: Find all the links and return them in an array")
    return []