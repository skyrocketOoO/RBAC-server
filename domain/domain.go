package domain

type Edge struct {
	UNs   string `json:"u_ns" bson:"u_ns"`
	UName string `json:"u_name" bson:"u_name"`
	Rel   string `json:"rel" bson:"rel"`
	VNs   string `json:"v_ns" bson:"v_ns"`
	VName string `json:"v_name" bson:"v_name"`
}

type Vertex struct {
	Ns   string `json:"ns"`
	Name string `json:"name"`
}

type TreeNode struct {
	Ns       string               `json:"ns"`
	Name     string               `json:"name"`
	Children map[string]*TreeNode `json:"children"`
}

type Permission struct {
	Rel  string
	Ns   string
	Name string
}
