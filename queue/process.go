package queue

import (
	"github.com/jinzhu/gorm"
	"github.com/nathan-osman/sprise/db"
)

// selectUpload attempts to retrieve an upload from the queue. If none was
// found, nil is returned.
func (q *Queue) selectUpload() (*db.Upload, error) {
	u := &db.Upload{}
	if err := q.conn.Preload("Bucket").Order("date").Find(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// createPhoto creates a Photo object and thumbnail from an upload.
func (q *Queue) createPhoto(name, thumbName string, u *db.Upload) (*db.Photo, error) {
	m, err := parseMetadata(name)
	if err != nil {
		return nil, err
	}
	t, err := generateThumbnail(name, thumbName, 200, 200)
	if err != nil {
		return nil, err
	}
	return &db.Photo{
		Date:     m.Date,
		Filename: u.Filename,
		Width:    t.originalWidth,
		Height:   t.originalHeight,
		Camera:   m.Camera,
		Bucket:   u.Bucket,
		BucketID: u.BucketID,
	}, nil
}

// processUpload processes an upload. This includes extracting the metadata
// from the image, generating a thumbnail, and uploading it to the intended
// bucket.
func (q *Queue) processUpload(u *db.Upload) error {
	q.log.Debugf("processing \"%s\"...", u.Filename)
	var (
		name      = u.Path(q.queueDir)
		thumbName = u.ThumbPath(q.queueDir)
	)
	p, err := q.createPhoto(name, thumbName, u)
	if err != nil {
		if err := q.conn.Delete(u).Error; err != nil {
			q.log.Error(err.Error())
		}
		return err
	}
	q.log.Debugf("uploading \"%s\"...", p.Filename)
	p.URL, err = uploadFile(name, p.Bucket)
	if err != nil {
		return err
	}
	q.log.Debugf("uploading \"%s\" thumbnail...", p.Filename)
	p.ThumbURL, err = uploadFile(thumbName, p.Bucket)
	if err != nil {
		return err
	}
	return q.conn.Transaction(func(conn *db.Conn) error {
		if err := conn.Delete(u).Error; err != nil {
			return err
		}
		return conn.Save(p).Error
	})
}
