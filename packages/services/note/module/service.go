package note

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
)

type Note struct {
	gorm.Model
	NoteId     uuid.UUID `json:"noteId"`
	UserId     uuid.UUID `json:"userId"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	Category   int       `json:"category"`
	TitleName  string    `json:"titleName"`
	DetailText string    `json:"detailText"`
}

func GetNotes(userId uuid.UUID, appId int, merchantId uuid.UUID, roleId int) map[string]interface{} {

	accs := &[]Note{}

	//check for errors and duplicate emails
	err := d.GetDB().Table("services.notes").Where("user_id = ? and app_id =? and merchant_id = ? and role_id = ?", userId, appId, merchantId, roleId).Scan(accs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11131:The user's note was not found for the specified application")
		}
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	response := u.Message(true, "0x11006:Requirement passed")
	response["notes"] = accs

	return response
}

func (note *Note) Create() map[string]interface{} {

	note.NoteId = uuid.New()

	d.GetDB().Table("services.notes").Create(note)

	if note.ID <= 0 {
		return u.Message(false, "0x11132:Failed to create note, connection error")
	}

	note.ID = 0

	response := u.Message(true, "0x11133:Note has been created")
	response["notes"] = note
	return response
}

func SetCategoryByNotId(noteId uuid.UUID, categoryId int, userId uuid.UUID, appId int, merchantId uuid.UUID, roleId int) map[string]interface{} {

	note, isOk := GetNote(noteId, userId, appId, merchantId, roleId)

	if !isOk {
		return u.Message(false, "0x11134:Unkown note error. Please retry")
	}

	note.Category = categoryId

	err := d.GetDB().Table("services.notes").Where("note_id = ? and user_id = ? and app_id =? and merchant_id = ? and role_id = ?", noteId, userId, appId, merchantId, roleId).Save(note).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	return u.Message(true, "0x11135:Note has been changed category")
}

func Delete(noteId uuid.UUID, userId uuid.UUID, appId int, merchantId uuid.UUID, roleId int) map[string]interface{} {

	note, isOk := GetNote(noteId, userId, appId, merchantId, roleId)

	if !isOk {
		return u.Message(false, "0x11134:Unkown note error. Please retry")
	}

	err := d.GetDB().Table("services.notes").Unscoped().Delete(&note).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	return u.Message(true, "0x11136:Note has been deleted")
}

func GetNote(noteId uuid.UUID, userId uuid.UUID, appId int, merchantId uuid.UUID, roleId int) (*Note, bool) {

	acc := &Note{}

	//check for errors and duplicate emails
	err := d.GetDB().Table("services.notes").Where("note_id = ? and user_id = ? and app_id =? and merchant_id = ? and role_id = ?", noteId, userId, appId, merchantId, roleId).First(acc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// return u.Message(false, "0x11004:Connection error. Please retry"), acc
		println(u.Message(false, "0x11004:Connection error. Please retry"))
		return acc, false
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			println(u.Message(false, "0x11131:The user's note was not found for the specified application"))
			return acc, false
		}
		println(u.Message(false, "0x11004:Connection error. Please retry"))
		return acc, false
	}

	println(u.Message(true, "0x11006:Requirement passed"))
	return acc, true
}
