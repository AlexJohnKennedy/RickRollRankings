package conditionchecker

import "fmt"
import "context"
import "github.com/bluele/gcache"

// LinkCacheSize is a constant which specifies the capacity of the link-checking LFU cache.
const LinkCacheSize = 100;

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
}
func buildNewLinkChecker(targetStrings []string, cacheSize int) *linkChecker {
	return &linkChecker{
		targetStrings: targetStrings,
		cacheSize: cacheSize,
		cache: gcache.New(cacheSize).LFU().Build(),
	};
}
func (lc *linkChecker) check(url string) bool {
	
}