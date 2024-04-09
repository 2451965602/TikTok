package db

import (
	"errors"
	"work4/bootstrap/env"
)

func CreateMessage(form, to, msg string) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	err := DB.
		Table(env.MsgTable).
		Create(&Message{
			FromUserId: form,
			ToUserId:   to,
			Content:    msg,
			Status:     0,
		}).
		Error

	if err != nil {
		return err
	}

	return nil
}

func GetMessage(to string) (*[]Message, error) {

	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var msglist []Message

	err := DB.
		Table(env.MsgTable).
		Where("status = ? ", 0).
		Where(`to_user_id = ?`, to).
		Find(&msglist).
		Error

	if err != nil {
		return nil, err
	}

	for _, msg := range msglist {
		DB.Table(env.MsgTable).
			Where(`msg_id = ?`, msg.MsgId).
			Update("status", 1)
	}

	return &msglist, nil
}
