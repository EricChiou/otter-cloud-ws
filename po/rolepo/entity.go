package rolepo

// Table name
const Table string = "role"

// PK name
const PK string = "id"

// column name
const (
	ID          string = "id"
	Code        string = "code"
	Name        string = "name"
	Lv          string = "lv"
	SortNo      string = "sort_no"
	Enable      string = "enable"
	CreatedDate string = "created_date"
	UpdatedDate string = "updated_date"
)

// Entity role table entity
type Entity struct {
	ID          int    `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Lv          int    `json:"lv,omitempty"`
	SortNo      int    `json:"sortNo,omitempty"`
	Enable      bool   `json:"enable,omitempty"`
	CreatedDate string `json:"creatDate,omitempty"`
	UpdatedDate string `json:"updateDate,omitempty"`
}
