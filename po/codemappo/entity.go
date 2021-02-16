package codemappo

// Table name
const Table string = "codemap"

// PK name
const PK string = "id"

// column name
const (
	ID          string = "id"
	Type        string = "type"
	Code        string = "code"
	Name        string = "name"
	SortNo      string = "sort_no"
	Enable      string = "enable"
	CreatedDate string = "created_date"
	UpdatedDate string = "updated_date"
)

// Entity codemap table entity
type Entity struct {
	ID          int    `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	SortNo      int    `json:"sortNo,omitempty"`
	Enable      bool   `json:"enable,omitempty"`
	CreatedDate string `json:"creatDate,omitempty"`
	UpdatedDate string `json:"updateDate,omitempty"`
}
