package models

// UserGroup : base model
type UserGroup struct {
	ID      int32 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}
