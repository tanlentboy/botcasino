package hashdb

// HashDatabase 数据库
type HashDatabase struct {
	bucketNum uint32    // 桶数量
	buckets   []*bucket // 桶列表
}

// Create 创建数据库
func Create(bucketNum uint32) (*HashDatabase, error) {
	if bucketNum == 0 {
		return nil, ErrCreatingDatabase
	}

	var i uint32
	buckets := make([]*bucket, 0, bucketNum)
	for i = 0; i < bucketNum; i++ {
		b := newBucket(i)
		buckets = append(buckets, b)
	}
	return &HashDatabase{buckets: buckets, bucketNum: bucketNum}, nil
}

// 获取所在桶
func (db *HashDatabase) getBucketID(key uint32) uint32 {
	return key % db.bucketNum
}

// Insert 插入数据
func (db *HashDatabase) Insert(key uint32, value interface{}) uint32 {
	bucketID := db.getBucketID(key)
	return db.buckets[bucketID].insert(key, value)
}

// Erase 擦除数据
func (db *HashDatabase) Erase(key uint32) uint32 {
	bucketID := db.getBucketID(key)
	return db.buckets[bucketID].erase(key)
}

// Update 更新数据
func (db *HashDatabase) Update(key uint32, fn func(interface{}) interface{}) {
	bucketID := db.getBucketID(key)
	db.buckets[bucketID].update(key, fn)
}

// Find 查找数据
func (db *HashDatabase) Find(key uint32) (interface{}, error) {
	bucketID := db.getBucketID(key)
	return db.buckets[bucketID].find(key)
}

// Clear 清空数据
func (db *HashDatabase) Clear() {
	var i uint32
	for i = 1; i <= db.bucketNum; i++ {
		db.buckets[i].clear()
	}
}
