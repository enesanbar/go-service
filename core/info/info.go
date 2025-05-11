package info

import (
	"encoding/json"
)

// Compile time variables
var (
	ServiceName      string
	ServiceNameHuman string
	Version          string
	CommitSHA        string
	BuildServer      string
	BuildDate        string
)

type BuildInfo struct {
	ServiceName      string `json:"service_name,omitempty"`
	ServiceNameHuman string `json:"service_name_human,omitempty"`
	Version          string `json:"version,omitempty"`
	CommitSHA        string `json:"commit_sha,omitempty"`
	BuildServer      string `json:"build_server,omitempty"`
	BuildDate        string `json:"build_date,omitempty"`
}

func (b BuildInfo) String() string {
	bytes, err := json.MarshalIndent(BuildInfo{
		ServiceName:      ServiceName,
		ServiceNameHuman: ServiceNameHuman,
		Version:          Version,
		CommitSHA:        CommitSHA,
		BuildServer:      BuildServer,
		BuildDate:        BuildDate,
	}, "", "   ")
	if err != nil {
		return "Unable to marshal struct"
	}

	return string(bytes)
}
