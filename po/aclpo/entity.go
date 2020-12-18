package aclpo

// table name
const Table string = "acl"

// pk name
const PK string = "id"

// column name
const (
	ID          string = "id"
	Code        string = "code"
	Name        string = "name"
	Type        string = "type"
	Lv          string = "lv"
	SortNo      string = "sort_no"
	Enable      string = "enable"
	CreatedDate string = "created_date"
	UpdatedDate string = "updated_date"
)

// Entity acl table entity
type Entity struct {
	ID          int    `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Lv          int    `json:"lv,omitempty"`
	SortNo      int    `json:"sortNo,omitempty"`
	Enable      bool   `json:"enable,omitempty"`
	CreatedDate string `json:"creatDate,omitempty"`
	UpdatedDate string `json:"updateDate,omitempty"`
}
