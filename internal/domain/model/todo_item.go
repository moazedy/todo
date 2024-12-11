package model

import (
	"time"

	"gorm.io/gorm"
)

// NOTE : as "D" principal of SOLID, abstraction must not has any dependency to technology or low/hig level modules.
// but for simplicity and also using generic repo layer, gorm tags has been applied in the service model.
type TodoItem struct {
	ID          string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileName    string    `json:"fileName"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ti TodoItem) IsEntity() {}
