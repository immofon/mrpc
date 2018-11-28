package rpc

type RetStatus string

const (
	Ok           RetStatus = "ok"
	Internal     RetStatus = "internal"
	Unauthorized RetStatus = "unauthorized"
	NotFound     RetStatus = "not_found"
)
