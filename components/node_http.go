package components

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

var partitionsRegex = regexp.MustCompile("/partition/([0-9a-zA-Z_\\-\\.]+)/([0-9a-zA-Z]+)$")

func (n Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method { // TODO make this accept only posts
		failRequest(w, "Invalid request for method")
		return
	}

	// extract the partition id
	matches := partitionsRegex.FindStringSubmatch(r.RequestURI)
	if len(matches) == 0 {
		failRequest(w, "Unknown endpoint: %s", r.RequestURI)
		return
	}

	// extract the payload
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		failRequest(w, "Invalid request body %v", err)
		return
	}
	defer r.Body.Close()

	topic := matches[1]
	partID := matches[2]
	log.Debug("Pushing data to partition: %s/%s", topic, partID)

	// now we should try and discover the topic - don't load anything we don't alreay have
	parts, exists := n.partitions[topic]
	if !exists {
		failRequest(w, "The topic %s does not exist on this node", topic)
		return
	}

	// do we have that actual partition?
	partIndex, err := strconv.Atoi(partID)
	if err != nil {
		failRequest(w, "The partition index %s is invalid", partID)
		return
	}

	// is it a valid partition?
	if partIndex >= len(parts) || partIndex < 0 {
		failRequest(w, "The partition index %d is invalid", partIndex)
		return
	}

	// all good in the hood - push it
	partition := parts[partIndex]
	partition.Write(string(bs))
	log.Debug("Successfully wrote to partition %s/%d", topic, partIndex)
}

func failRequest(w http.ResponseWriter, errFmt string, args ...interface{}) {
	errMsg := fmt.Sprintf(errFmt, args)
	log.Warn(errMsg)
	http.Error(w, errMsg, http.StatusBadRequest)
}
