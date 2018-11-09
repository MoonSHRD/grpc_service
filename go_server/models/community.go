package models

type Community struct {
    Id     string `bson:"_id,omitempty" json:"id"`
    Name   string `json:"name"`
    Avatar string `json:"avatar"`
    Type   string `json:"type"`
    Admin  string `json:"admin"`
}
