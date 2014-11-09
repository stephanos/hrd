package internal

import (
	"fmt"

	"appengine/datastore"
)

// keyStringID returns the ID of the passed-in Key as a string.
func keyStringID(key *datastore.Key) (id string) {
	id = key.StringID()
	if id == "" && key.IntID() > 0 {
		id = fmt.Sprintf("%v", key.IntID())
	}
	return
}

// KeyString returns a string representation of the passed-in Key.
func KeyString(key *datastore.Key) string {
	if key == nil {
		return ""
	}
	ret := fmt.Sprintf("Key{'%v', %v}", key.Kind(), keyStringID(key))
	parent := KeyString(key.Parent())
	if parent != "" {
		ret += fmt.Sprintf("[Parent%v]", parent)
	}
	return ret
}
