package kafka

type PosDTO struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Z  int     `json:"z"`
	A  int     `json:"a"`
	S  int     `json:"s"`
	St int     `json:"st"`
}

type MessageDTO struct {
	T      int64                  `json:"t"`
	ST     int                    `json:"st"`
	Pos    PosDTO                 `json:"pos"`
	Params map[string]interface{} `json:"p"`
}
