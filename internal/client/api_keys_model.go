package client

type ApiKeyRequest struct {
	Alias  string        `json:"alias"`
	Status string        `json:"status,omitempty"`
	Scopes *ApiKeyScopes `json:"scopes,omitempty"`
}

type ApiKeyScopes struct {
	ApiClients    *ApiKeyScope `json:"api_clients,omitempty"`
	Channels      *ApiKeyScope `json:"channels,omitempty"`
	Domains       *ApiKeyScope `json:"domains,omitempty"`
	Events        *ApiKeyScope `json:"events,omitempty"`
	Listeners     *ApiKeyScope `json:"listeners,omitempty"`
	Messages      *ApiKeyScope `json:"messages,omitempty"`
	Queues        *ApiKeyScope `json:"queues,omitempty"`
	Schemas       *ApiKeyScope `json:"schemas,omitempty"`
	Stats         *ApiKeyScope `json:"stats,omitempty"`
	Subscriptions *ApiKeyScope `json:"subscriptions,omitempty"`
	Topics        *ApiKeyScope `json:"topics,omitempty"`
}

type ApiKeyScope struct {
	Access  string   `json:"access,omitempty"`
	Targets []string `json:"targets,omitempty"`
}

type ApiKeyResponse struct {
	Id         string        `json:"id"`
	Alias      string        `json:"alias"`
	Status     string        `json:"status"`
	Key        string        `json:"key"`
	Scopes     *ApiKeyScopes `json:"scopes,omitempty"`
	CreatedAt  int64         `json:"created_at,omitempty"`
	UpdatedAt  int64         `json:"expires_at,omitempty"`
	LastUsedAt int64         `json:"last_used_at,omitempty"`
}
