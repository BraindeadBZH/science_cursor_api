package api

import (
	"crypto/rand"
	"fmt"
	"time"
)

func sessionExists(key []byte) (bool, error) {
	count, err := dbHandle.Model(&sessionModel{}).Where("key = ?", key).Count()
	if err != nil {
		return false, fmt.Errorf("Could not check session existance: '%s'", err.Error())
	}
	return count > 0, nil
}

func getSession(key []byte) (*sessionModel, error) {
	session := &sessionModel{}
	err := dbHandle.Model(&sessionModel{}).Where("key = ?", key).Select(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func createSession() (*sessionModel, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("Could not create session key: '%s'", err.Error())
	}

	session := &sessionModel{}
	session.Key = key
	session.Data = sessionData{Auth: false, MarkedForDeletion: false}
	session.LastUsed = time.Now()

	err = dbHandle.Insert(session)
	if err != nil {
		return nil, fmt.Errorf("Error while creating session: '%s'", err.Error())
	}

	return session, nil
}

func (session *sessionModel) save() error {
	session.LastUsed = time.Now()
	return dbHandle.Update(session)
}

func (session *sessionModel) destroy() error {
	return dbHandle.Delete(session)
}
