package model

type Route struct {
	Id         string `json:"id"`
	Method     string `json:"method"`
	Pattern    string `json:"url_pattern"`
	Entrypoint string `json:"entrypoint"`
	Command    string `json:"command"`
}
