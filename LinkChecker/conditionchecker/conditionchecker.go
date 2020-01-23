package conditionchecker

import "context"
import "github.com/bluele/gcache"
import "net/http"
import u "net/url"
import "strings"
import "errors"

// LinkCacheSize is a constant which specifies the capacity of the link-checking LFU cache.
const LinkCacheSize = 100;
const maxRedirects = 10;	

// Message is a simple struct containing the data of a incoming 'links to check' kafka message.
type Message struct {
	originalData string
	linksToCheck []string
}

// LaunchMessageChecker is a function which waits for messages to on an input channel, and routes those messages to either a 'match' or 'no match'
// output channel, based on where the re-directs go on the links to check in the message. This will launch child threads to handle link examination,
// and cache results in a LFU cache to avoid re-fetching information from the same link over and over again.
// It is expected that this function is launched as a go-routine.
func LaunchMessageChecker(ctx context.Context, targetStrings []string, inputs chan *Message, matches chan *Message, nonmatches chan *Message) {
	checker := buildNewLinkChecker(targetStrings, LinkCacheSize);
}

type linkChecker struct {
	targetStrings []string
	cacheSize int			// The specified maximum size of the link cache
	cache gcache.Cache		// A Least-frequently-used cache to remember results of links we have already looked at. This is apparently thread-safe.
	httpClient *http.Client	// The standard-lib http client we will use to probe links. This is also thread-safe according to godocs.
}
func buildNewLinkChecker(targetStrings []string, cacheSize int) *linkChecker {
	// Setup our custom redirect-handler function for the httpClient, which is doing the actual 'checking' logic.
	// We are leveraging the fact that the net-package http client can be configured such that this custom function intercepts
	// each redirect request just before it is followed by the client. This way, our function can 'inspect' the redirect-url to determine
	// if it is targetting a rick-roll url. If this ever occurs, we will abort the request by returning an error such that the link checker
	// realises the abort was due to a positive match. Any other error type will result in a negative result.
	customRedirectHandler := func(req *http.Request, via []*http.Request) error {
		if len(via) == maxRedirects {
			return &tooManyRedirectsError{ "Too many redirects! Max limit reached, probe aborted" };
		}
		for _, targ := range targetStrings {
			if strings.Contains(req.RequestURI, targ) {
				return &foundMatchPseudoError{ req.RequestURI };
			}
		}
		return nil;
	};
	
	return &linkChecker{
		targetStrings: targetStrings,
		cacheSize: cacheSize,
		cache: gcache.New(cacheSize).LFU().Build(),
		httpClient: &http.Client{
			CheckRedirect: customRedirectHandler,
			Timeout: 10000000000,	// Ten seconds
		},
	};
}
func (lc *linkChecker) check(url string) bool {
	// If we have checked this url before, then we can return instantly
	if (lc.cache.Has(url)) {
		r, _ := lc.cache.Get(url);
		return r.(bool);
	}
	// We have not checked this url before, so we must probe it with the httpclient
	resp, err := lc.httpClient.Get(url);
	if (err != nil) {
		urlErr := err.(*u.Error).Unwrap();
		if _, ok := urlErr.(*foundMatchPseudoError); ok {
			lc.cache.Set(url, true);
			return true;
		}
		lc.cache.Set(url, false);
		return false;
	}
	defer resp.Body.Close();
	lc.cache.Set(url, false);
	return false;
}

// Custom Error implementations, used as signals to the link checker object
type tooManyRedirectsError struct {
	msg string
}
func (e *tooManyRedirectsError) Error() string {
	return e.msg;
}

type foundMatchPseudoError struct {
	matchingURL string
}
func (e *foundMatchPseudoError) Error() string {
	return e.matchingURL;
}