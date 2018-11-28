package mrpc

type RetStatus string

const (
	Ok           RetStatus = "ok"
	Internal     RetStatus = "internal"
	Network      RetStatus = "network"
	Unauthorized RetStatus = "unauthorized"
	NotFound     RetStatus = "not_found"
)
