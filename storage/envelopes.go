package storage

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"

	"github.com/boltdb/bolt"
)

var (
	// ErrNothingLeft 领完了
	ErrNothingLeft = errors.New("nothing left")
	// ErrRepeatReceive 重复领取
	ErrRepeatReceive = errors.New("repeat receive")
	// ErrNotActivated 没有激活
	ErrNotActivated = errors.New("not activated")
	// ErrAlreadyActivated 已经激活
	ErrAlreadyActivated = errors.New("already activated")
	// ErrPermissionDenied 没有权限
	ErrPermissionDenied = errors.New("permission denied")
	// ErrRedEnvelopedExpired 红包已过期
	ErrRedEnvelopedExpired = errors.New("red envelope expired")
)

// ********************** 结构图 **********************
// {
// 	"red_envelopes": {
// 		"sid": {
// 			"seq": 0,					// 红包领取序列
// 			"min": 0,					// 手气最烂序列
// 			"max": 0,					// 手气最佳序列
// 			"base": types.RedEnvelope	// 红包基本信息
//			"users": {					// 红包已领用户
//				"user_id": ""
//			}
// 			"record": {					// 红包领取记录
// 				"seq": types.RedEnvelopeRecord
// 			}
//			"expired": true				// 红包过期
// 		},
//		sequeue: 0						// 红包ID生成序列
//		last_expired: 0,				// 上次过期红包ID
// 	}
// }
// ***************************************************

// RedEnvelopeStorage 红包存储
type RedEnvelopeStorage struct {
}

// 创建领取记录
func (handler *RedEnvelopeStorage) newRecord(tx *bolt.Tx, sid string,
	envelopes []int) (int, int, error) {

	minValueSeq, maxValueSeq := 0, 0
	minValue, maxValue := math.MaxInt32, 0
	recordBucket, err := ensureBucketExists(tx, "red_envelopes", sid, "record")
	if err != nil {
		return 0, 0, err
	}
	for i := range envelopes {
		seq, err := recordBucket.NextSequence()
		if err != nil {
			return 0, 0, err
		}

		val := RedEnvelopeRecord{Value: envelopes[i]}
		jsb, err := json.Marshal(&val)
		if err != nil {
			return 0, 0, err
		}

		sseq := strconv.FormatUint(seq, 10)
		err = recordBucket.Put([]byte(sseq), jsb)
		if err != nil {
			return 0, 0, err
		}

		if envelopes[i] < minValue {
			minValue = envelopes[i]
			minValueSeq = int(seq)
		}

		if envelopes[i] > maxValue {
			maxValue = envelopes[i]
			maxValueSeq = int(seq)
		}
	}
	return minValueSeq, maxValueSeq, nil
}

// NewRedEnvelope 创建新红包
func (handler *RedEnvelopeStorage) NewRedEnvelope(data *RedEnvelope, envelopes []int) (*RedEnvelope, error) {
	err := blotDB.Update(func(tx *bolt.Tx) error {
		// 生成红包ID
		rootBucket, err := ensureBucketExists(tx, "red_envelopes")
		if err != nil {
			return err
		}
		if rootBucket.Sequence() < DefaultRenvelopesID {
			if err = rootBucket.SetSequence(DefaultRenvelopesID); err != nil {
				return err
			}
		}
		data.ID, err = rootBucket.NextSequence()
		if err != nil {
			return err
		}

		// 序列化数据
		data.Received = 0
		data.Active = false
		jsb, err := json.Marshal(data)
		if err != nil {
			return err
		}

		// 插入基本信息
		sid := strconv.FormatUint(data.ID, 10)
		bucket, err := ensureBucketExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("base"), jsb)
		if err != nil {
			return err
		}

		// 插入领取用户
		_, err = ensureBucketExists(tx, "red_envelopes", sid, "users")
		if err != nil {
			return err
		}

		// 插入领取记录
		minValueSeq, maxValueSeq, err := handler.newRecord(tx, sid, envelopes)
		if err != nil {
			return err
		}

		// 插入已领取序列
		err = bucket.Put([]byte("seq"), []byte("0"))
		if err != nil {
			return err
		}

		// 插入手气最烂序列
		err = bucket.Put([]byte("min"), []byte(strconv.Itoa(minValueSeq)))
		if err != nil {
			return err
		}

		// 插入手气最佳序列
		err = bucket.Put([]byte("max"), []byte(strconv.Itoa(maxValueSeq)))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}

// IsExpired 是否过期
func (handler *RedEnvelopeStorage) IsExpired(id uint64) bool {
	var expired bool
	sid := strconv.FormatUint(id, 10)
	err := blotDB.View(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}
		expired = bucket.Get([]byte("expired")) != nil
		return nil
	})

	if err != nil {
		return false
	}
	return expired
}

// SetExpired 设置过期
func (handler *RedEnvelopeStorage) SetExpired(id uint64) error {
	sid := strconv.FormatUint(id, 10)
	return blotDB.Update(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}
		return bucket.Put([]byte("expired"), []byte("true"))
	})
}

// IsReceived 是否已领取
func (handler *RedEnvelopeStorage) IsReceived(id uint64, userID int32) (bool, error) {
	received := false
	sid := strconv.FormatUint(id, 10)
	err := blotDB.View(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid, "users")
		if err != nil {
			return err
		}
		received = bucket.Get([]byte(strconv.FormatInt(int64(userID), 10))) != nil
		return nil
	})
	if err != nil {
		return received, err
	}
	return received, nil
}

// GetLastExpired 获取上次过期红包
func (handler *RedEnvelopeStorage) GetLastExpired() (uint64, error) {
	var id uint64
	err := blotDB.View(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes")
		if err != nil {
			return err
		}

		sid := bucket.Get([]byte("last_expired"))
		if sid == nil {
			return nil
		}

		id, err = strconv.ParseUint(string(sid), 10, 64)
		if err != nil {
			return nil
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}

// SetLastExpired 设置上次过期红包
func (handler *RedEnvelopeStorage) SetLastExpired(id uint64) error {
	sid := strconv.FormatUint(id, 10)
	return blotDB.Update(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes")
		if err != nil {
			return err
		}
		return bucket.Put([]byte("last_expired"), []byte(sid))
	})
}

// GetRedEnvelope 获取红包信息
func (handler *RedEnvelopeStorage) GetRedEnvelope(id uint64) (*RedEnvelope, uint32, error) {
	var received uint32
	var base RedEnvelope
	sid := strconv.FormatUint(id, 10)
	err := blotDB.View(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}

		// 获取红包信息
		jsb := bucket.Get([]byte("base"))
		if err = json.Unmarshal(jsb, &base); err != nil {
			return err
		}

		// 已领取数量
		seq := bucket.Get([]byte("seq"))
		numReceived, err := strconv.Atoi(string(seq))
		if err != nil {
			return err
		}

		// 剩余红包数量
		received = uint32(numReceived)
		return nil
	})

	if err != nil {
		return nil, 0, err
	}
	return &base, received, nil
}

// ActiveRedEnvelope 激活红包
func (handler *RedEnvelopeStorage) ActiveRedEnvelope(id uint64, userID, chatID int64, messageID int32) error {
	sid := strconv.FormatUint(id, 10)
	return blotDB.Update(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}

		// 检查状态
		if bucket.Get([]byte("expired")) != nil {
			return ErrRedEnvelopedExpired
		}

		// 获取红包基本信息
		var base RedEnvelope
		jsb := bucket.Get([]byte("base"))
		if err = json.Unmarshal(jsb, &base); err != nil {
			return err
		}

		if base.Active {
			return ErrAlreadyActivated
		}

		if base.SenderID != userID {
			return ErrPermissionDenied
		}

		// 更新红包基本信息
		base.Active = true
		base.GroupID = chatID
		base.MessageID = messageID
		if jsb, err = json.Marshal(&base); err != nil {
			return err
		}
		return bucket.Put([]byte("base"), jsb)
	})
}

// 领取红包
func (handler *RedEnvelopeStorage) receiveRedEnvelope(tx *bolt.Tx, sid string, seq int,
	user *RedEnvelopeUser) (int, error) {

	recordBucket, err := getBucketIfExists(tx, "red_envelopes", sid, "record")
	if err != nil {
		return 0, err
	}

	var record RedEnvelopeRecord
	key := []byte(strconv.Itoa(seq))
	jsb := recordBucket.Get(key)
	if err = json.Unmarshal(jsb, &record); err != nil {
		return 0, err
	}
	record.User = user

	jsb, err = json.Marshal(&record)
	if err != nil {
		return 0, err
	}

	if err = recordBucket.Put(key, jsb); err != nil {
		return 0, err
	}
	return record.Value, nil
}

// ReceiveRedEnvelope 领取红包
func (handler *RedEnvelopeStorage) ReceiveRedEnvelope(id uint64, userID int32,
	firstName string) (int, int, error) {

	received, err := handler.IsReceived(id, userID)
	if err != nil {
		return 0, 0, err
	}

	if received {
		return 0, 0, ErrRepeatReceive
	}

	value := 0
	count := 0
	sid := strconv.FormatUint(id, 10)
	err = blotDB.Update(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}

		// 检查状态
		if bucket.Get([]byte("expired")) != nil {
			return ErrRedEnvelopedExpired
		}

		// 已领取数量
		seq := bucket.Get([]byte("seq"))
		numReceived, err := strconv.Atoi(string(seq))
		if err != nil {
			return err
		}

		// 红包是否充足
		var base RedEnvelope
		jsb := bucket.Get([]byte("base"))
		if err = json.Unmarshal(jsb, &base); err != nil {
			return err
		}
		if uint32(numReceived) >= base.Number {
			return ErrNothingLeft
		}

		// 红包是否激活
		if !base.Active {
			return ErrNotActivated
		}

		// 是否重复领取
		usersBucket, err := getBucketIfExists(tx, "red_envelopes", sid, "users")
		if err != nil {
			return err
		}
		key := []byte(strconv.FormatInt(int64(userID), 10))
		if usersBucket.Get(key) != nil {
			return ErrRepeatReceive
		}

		// 执行领取红包
		newSeq := numReceived + 1
		value, err = handler.receiveRedEnvelope(tx, sid, newSeq, &RedEnvelopeUser{
			UserID:    userID,
			FirstName: firstName,
		})
		if err != nil {
			return err
		}
		base.Received += uint32(value)

		// 更新红包信息
		if jsb, err = json.Marshal(&base); err != nil {
			return err
		}
		if err = bucket.Put([]byte("base"), jsb); err != nil {
			return err
		}
		if err = usersBucket.Put(key, []byte("")); err != nil {
			return err
		}
		if err = bucket.Put([]byte("seq"), []byte(strconv.Itoa(newSeq))); err != nil {
			return err
		}

		count = int(base.Number - uint32(newSeq))
		return nil
	})

	if err != nil {
		return 0, 0, err
	}
	return value, count, nil
}

// GetTwoTxtremes 获取两个极端
func (handler *RedEnvelopeStorage) GetTwoTxtremes(id uint64) (*RedEnvelopeRecord, *RedEnvelopeRecord, error) {
	var minRecord RedEnvelopeRecord
	var maxRecord RedEnvelopeRecord
	sid := strconv.FormatUint(id, 10)
	err := blotDB.View(func(tx *bolt.Tx) error {
		bucket, err := getBucketIfExists(tx, "red_envelopes", sid)
		if err != nil {
			return err
		}

		// 获取极端情况序列号
		minValueSeq := bucket.Get([]byte("min"))
		maxValueSeq := bucket.Get([]byte("max"))
		if minValueSeq == nil || maxValueSeq == nil {
			return errors.New("nou found")
		}

		// 获取红包序列号信息
		recordBucket, err := getBucketIfExists(tx, "red_envelopes", sid, "record")
		if err != nil {
			return err
		}

		minValueData := recordBucket.Get(minValueSeq)
		maxValueData := recordBucket.Get(maxValueSeq)
		if minValueData == nil || maxValueData == nil {
			return errors.New("nou found")
		}

		if err = json.Unmarshal(minValueData, &minRecord); err != nil {
			return err
		}
		if err = json.Unmarshal(maxValueData, &maxRecord); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return &minRecord, &maxRecord, nil
}

// ForeachRedEnvelopes 遍历红包列表
func (handler *RedEnvelopeStorage) ForeachRedEnvelopes(startID uint64, callback func(*RedEnvelope)) error {
	var base RedEnvelope
	return blotDB.View(func(tx *bolt.Tx) error {
		rootBucket, err := getBucketIfExists(tx, "red_envelopes")
		if err != nil {
			return err
		}

		cursor := rootBucket.Cursor()
		seek := []byte(strconv.FormatUint(startID, 10))
		for k, v := cursor.Seek(seek); k != nil && v == nil; k, v = cursor.Next() {
			if bucket := rootBucket.Bucket(k); bucket != nil {
				jsb := bucket.Get([]byte("base"))
				if err = json.Unmarshal(jsb, &base); err != nil {
					continue
				}
				if callback != nil {
					callback(&base)
				}
			}
		}
		return nil
	})
}
