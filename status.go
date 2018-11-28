package mrpc

type RetStatus string

const (
	Ok           RetStatus = "ok"           // success
	Network      RetStatus = "network"      // anything about network
	Internal     RetStatus = "internal"     // most error
	Unauthorized RetStatus = "unauthorized" // anything about auth
	NotFound     RetStatus = "not_found"    // not found something required
)
