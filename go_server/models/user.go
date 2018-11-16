package models

type User struct {
    Id          string `bson:"_id,omitempty" json:"id"`
    Firstname   string `json:"firstname"`
    Lastname    string `json:"lastname"`
    Bio         string `json:"bio"`
    Avatar      string `json:"avatar"`
    Last_active int    `json:"last_active"`
}

//type UserNonce struct {
//    Id    string `bson:"_id,omitempty"`
//    Nonce string
//}

