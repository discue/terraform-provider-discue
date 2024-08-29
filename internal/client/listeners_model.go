package client

type Listener struct {
	Id          string `json:"id,omitempty"`
	Alias       string `json:"alias"`
	Status      string `json:"status,omitempty"`
	NotifyUrl   string `json:"notify_url,omitempty"`
	LivenessUrl string `json:"liveness_url,omitempty"`
}

type ListenerRequest = Listener
type ListenerResponse = Listener
