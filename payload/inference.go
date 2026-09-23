package payload

type ExecuteInferenceResponse struct {
	Result string  `json:"result"`
	Live   float32 `json:"live"`
	Spoof  float32 `json:"spoof"`
}
