package model

type User struct {
	Id          uint64 `json:"id"`
	Username    string `json:"user_name"`
	CreatedTime uint32 `json:"created_time"`
	UpdatedTime uint32 `json:"updated_time"`
}
