package storage

import (
	"context"
	"fmt"
	"os"
    "strings"
	"time"

	"github.com/vinnymaker18/sms-on-pc/backend/common"
)

const (
	newSmsFetchQuery = "SELECT msg_time, origin, body, sms_id FROM sms WHERE user_id = $1 AND NOT seen"

	newSmsStoreQuery = "INSERT INTO sms (user_id, msg_time, origin, body) VALUES ($1, $2, $3, $4)"

	deleteOldSmsQuery = "DELETE FROM sms WHERE msg_time < $1"

	markSmsQuery = "UPDATE sms SET seen = 't' WHERE sms_id in $1"
)

// StoreNewSMS persists a new SMS in the database.
func StoreNewSMS(message *common.SMSMessage) error {
	conn, err := newDBConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database "+err.Error())
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), newSmsStoreQuery, message.UserID, message.Time,
		message.OriginAddress, message.MsgBody)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	return nil
}

// FetchNewSMS retrieves new text messages for the given userID.
func FetchNewSMS(userID int64) []*common.SMSMessage {
	conn, err := newDBConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database "+err.Error())
		return []*common.SMSMessage{}
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), newSmsFetchQuery, userID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error querying database"+err.Error())
		return []*common.SMSMessage{}
	}

	messages := make([]*common.SMSMessage, 0)

	for rows.Next() {
		var msgTime time.Time
		var originAddress string
		var msgBody string
		var msgID int64
		err = rows.Scan(&msgTime, &originAddress, &msgBody, &msgID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing returned database row"+err.Error())
			return []*common.SMSMessage{}
		}

		messages = append(messages, &common.SMSMessage{
			UserID:        userID,
			Time:          msgTime,
			MsgBody:       msgBody,
			OriginAddress: originAddress,
			MsgID:         msgID,
		})
	}

	return messages
}

func listifyIDs(ids []int64) string {
    var builder strings.Builder
    builder.WriteString("(")
    for _, id := range ids {
        builder.WriteString(fmt.Sprintf("%d,", id))
    }
    builder.WriteString("-1)")
    return builder.String()
}

// MarkAsRead marks the given message ids as read in the database.
func MarkAsRead(msgIDs []int64) error {
	conn, err := newDBConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database"+err.Error())
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), markSmsQuery, listifyIDs(msgIDs))
	return err
}

// DeleteOldSMS deletes old sms messages from the database.
func DeleteOldSMS() error {
	conn, err := newDBConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database "+err.Error())
		return err
	}
	defer conn.Release()

	cutOffTime := time.Now().AddDate(0, 0, -2)
	_, err = conn.Query(context.Background(), deleteOldSmsQuery, cutOffTime)
	return nil
}
