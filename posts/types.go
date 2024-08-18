package posts

type SendPostStrategy string

const (
	SendPostStratSeq  SendPostStrategy = "sequential"
	SendPostStratASeq SendPostStrategy = "async_sequential"
	SendPostStratBulk SendPostStrategy = "bulk"
	SendPostStratMass SendPostStrategy = "mass"
)
