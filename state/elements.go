// Package state provides elements that will be stored in the distributed
// key value store. This wile will contain the structures that make up the
// values
package state

type ClusterAdminCmd struct {
}

type SessionAdminCmd struct {
}

type SessionState struct {
	SessionId string
}
