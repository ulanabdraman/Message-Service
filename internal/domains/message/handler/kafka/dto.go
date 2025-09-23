package kafka

type PosDTO struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Z  int     `json:"z"`
	A  int     `json:"a"`
	S  int     `json:"s"`
	Sl int     `json:"sl"`
}

type MessageDTO struct {
	ID     int64                  `json:"id"`
	DT     int64                  `json:"dt"`
	ST     int64                  `json:"st"`
	Pos    PosDTO                 `json:"pos"`
	Params map[string]interface{} `json:"p"`
}
