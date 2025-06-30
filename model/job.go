// package model
//
// import "time"
//
//	type Job struct {
//	   ID        int       `db:"id" json:"id"`
//	   Payload   string    `db:"payload" json:"payload"`
//	   Status    string    `db:"status" json:"status"`
//	   Result    string    `db:"result" json:"result,omitempty"`
//	   CreatedAt time.Time `db:"created_at" json:"created_at"`
//	   UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
//	}
package model

import "time"

type Job struct {
	ID        int       `db:"id" json:"id"`
	Payload   string    `db:"payload" json:"payload"`
	Status    string    `db:"status" json:"status"`
	Result    *string   `db:"result" json:"result"` // Nullable
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
