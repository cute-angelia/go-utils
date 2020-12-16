package store

//统一化参数结构体
type ImageParam struct {
	Name    string
	Info    string
	Type    string
	Content *[]byte
}

//统一化返回结构体
type ImageReturn struct {
	Url    string
	Delete string
	Path   string
	ID     int
	Other  interface{}
}

//Proxy
type ProxyConf struct {
	Status bool   `json:"status"`
	Node   string `json:"node"`
}

// GitHub
type GithubConfig struct {
	Token  string // 开发者 token
	Owner  string // owner
	Repo   string // repo
	Branch string // branch
}

