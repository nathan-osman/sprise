package queue

import (
	"github.com/jinzhu/gorm"
	"github.com/nathan-osman/sprise/db"
)

// selectUpload attempts to retrieve an upload from the queue. If none was
// found, nil is returned.
func (q *Queue) selectUpload() (*db.Upload, error) {
	u := &db.Upload{}
	if err := q.conn.Order("date").Find(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// processUpload processes the upload.
func (q *Queue) processUpload(u *db.Upload) error {
	//...

	return nil
}
