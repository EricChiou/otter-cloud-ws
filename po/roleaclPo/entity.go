package roleaclpo

// table name
const Table string = "role_acl"

// pk name
const PK string = "id"

// column name
const (
	ID          string = "id"
	RoleCode    string = "role_code"
	ACLCode     string = "acl_code"
	CreatedDate string = "created_date"
	UpdatedDate string = "updated_date"
)

// Entity role_acl table entity
type Entity struct {
	ID          int    `json:"id,omitempty"`
	RoleCode    string `json:"roleCode,omitempty"`
	ACLCode     string `json:"aclCode,omitempty"`
	CreatedDate string `json:"creatDate,omitempty"`
	UpdatedDate string `json:"updateDate,omitempty"`
}
