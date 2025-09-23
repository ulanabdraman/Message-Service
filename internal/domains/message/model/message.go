package model

type Pos struct {
	X  float64 `json:"x" bson:"x"`
	Y  float64 `json:"y" bson:"y"`
	Z  int     `json:"z" bson:"z"`
	A  int     `json:"a" bson:"a"`
	S  int     `json:"s" bson:"s"`
	Sl int     `json:"sl" bson:"sl"`
}

type Message struct {
	ID     int64                  `json:"id"`
	DT     int64                  `json:"dt" bson:"dt"`
	ST     int64                  `json:"st" bson:"st"`
	Pos    Pos                    `json:"pos" bson:"pos"`
	Params map[string]interface{} `json:"p" bson:"p"`
}
