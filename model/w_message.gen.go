// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameWMessage = "w_message"

// WMessage mapped from table <w_message>
type WMessage struct {
	ID            int32  `gorm:"column:id;primaryKey" json:"id"`
	Messaage      string `gorm:"column:messaage" json:"messaage"`
	FromUsername  string `gorm:"column:from_username" json:"from_username"`
	BotName       string `gorm:"column:bot_name" json:"bot_name"`
	FromFirstName string `gorm:"column:from_first_name" json:"from_first_name"`
	FromLastName  string `gorm:"column:from_last_name" json:"from_last_name"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	ChatID        int64  `gorm:"column:chat_id" json:"chat_id"`
}

// TableName WMessage's table name
func (*WMessage) TableName() string {
	return TableNameWMessage
}