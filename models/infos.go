package models

type Infos struct {
	Infos []Info `json:"versoes"`
}

type Info struct {
	MessageID  string
	ID         string `json:"id"`
	Format     string `json:"versao"`
	Size       string `json:"tamanho"`
	Resolution string `json:"resolucao"`
}
