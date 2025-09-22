package model

type Pos struct {
	X  float64 `json:"x" bson:"x"`
	Y  float64 `json:"y" bson:"y"`
	Z  int     `json:"z" bson:"z"`
	A  int     `json:"a" bson:"a"`
	S  int     `json:"s" bson:"s"`
	St int     `json:"st" bson:"st"`
}

type Message struct {
	UUID   int64                  `json:"uuid" bson:"uuid"`
	T      int64                  `json:"t" bson:"t"`
	ST     int                    `json:"st" bson:"st"`
	Pos    Pos                    `json:"pos" bson:"pos"`
	Params map[string]interface{} `json:"p" bson:"p"`
}
