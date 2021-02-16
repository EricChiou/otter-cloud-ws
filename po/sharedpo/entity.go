package sharedpo

// Table name
const Table string = "shared_folder"

// PK name
const PK string = "id"

// column name
const (
	ID          string = "id"
	OwnerAcc    string = "owner_acc"
	SharedAcc   string = "shared_acc"
	BucketName  string = "bucket_name"
	Prefix      string = "prefix"
	Permission  string = "permission"
	CreatedDate string = "created_date"
	UpdatedDate string = "updated_date"
)

// Entity shared_folder table entity
type Entity struct {
	ID          int    `json:"id,omitempty"`
	OwnerAcc    string `json:"ownerAcc,omitempty"`
	SharedAcc   string `json:"sharedAcc,omitempty"`
	BucketName  string `json:"bucketName,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Permission  string `json:"permission,omitempty"`
	CreatedDate string `json:"creatDate,omitempty"`
	UpdatedDate string `json:"updateDate,omitempty"`
}
