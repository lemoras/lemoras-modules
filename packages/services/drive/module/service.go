package drive

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
)

type Bucket struct {
	gorm.Model
	BucketId   uuid.UUID `json:"bucketId"`
	UserId     uuid.UUID `json:"userId"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	BucketName string    `json:"bucketName"`
	EmptySize  float64   `json:"emptySize"`
	SizeLimit  float64   `json:"sizeLimit"`
}

type BucketItem struct {
	gorm.Model
	BucketItemId uuid.UUID `json:"bucketItemId"`
	BucketId     uuid.UUID `json:"bucketId"`
	ParentId     uuid.UUID `json:"parentId"`
	ItemType     int       `json:"itemId"`
	ItemName     string    `json:"itemName"`
	ItemSize     float64   `json:"itemSize"`
	ItemURL      string    `json:"itemURL"`
}

type Response struct {
	BucketId    uuid.UUID      `json:"bucketId"`
	BucketName  string         `json:"bucketName"`
	EmptySize   float64        `json:"emptySize"`
	SizeLimit   float64        `json:"sizeLimit"`
	BucketItems []ResponceItem `json:"bucketItems"`
}

// ResponceItem modeli (recursive)
type ResponceItem struct {
	BucketItemId uuid.UUID      `json:"bucketItemId"`
	ParentId     uuid.UUID      `json:"parentId"`
	ItemType     int            `json:"itemType"`
	ItemName     string         `json:"itemName"`
	ItemSize     float64        `json:"itemSize"`
	ItemURL      string         `json:"itemURL"`
	ItemCount    int            `json:"itemCount"`
	TotalSize    float64        `json:"totalSize"`
	BucketItems  []ResponceItem `json:"bucketItems,omitempty"`
}

type flatResult struct {
	BucketItemId uuid.UUID  `gorm:"column:bucket_item_id"`
	BucketId     uuid.UUID  `gorm:"column:bucket_id"`
	BucketName   string     `gorm:"column:bucket_name"`
	EmptySize    float64    `gorm:"column:empty_size"`
	SizeLimit    float64    `gorm:"column:size_limit"`
	ItemName     string     `gorm:"column:item_name"`
	ItemType     int        `gorm:"column:item_type"`
	ItemSize     float64    `gorm:"column:item_size"`
	ItemURL      string     `gorm:"column:item_url"`
	ParentId     *uuid.UUID `gorm:"column:parent_id"`
	Depth        int        `gorm:"column:depth"`
}

// GetBucketsWithItems veritabanından recursive bucket ve itemları getirir
func GetBucketsWithItems(
	userId, merchantId uuid.UUID,
	appId, roleId int,
	bucketId, bucketItemId uuid.UUID, // opsiyonel filtreler
	maxDepth int,
	initCompleted bool,
) map[string]interface{} {

	if !initCompleted {
		var existingBucketCount int64
		err2 := d.GetDB().Table("services.buckets").
			Where("user_id = ? AND merchant_id = ? AND app_id = ? AND role_id = ?", userId, merchantId, appId, roleId).
			Count(&existingBucketCount).Error

		if err2 != nil {
			return u.Message(false, "0x11004:Connection error (count). Please retry")
		}

		if existingBucketCount == 0 {
			newBucket := Bucket{
				BucketId:   uuid.New(),
				UserId:     userId,
				RoleId:     roleId,
				AppId:      appId,
				MerchantId: merchantId,
				BucketName: "Default Bucket",
				EmptySize:  500 * 1024 * 1024, // 500 MB in bytes
				SizeLimit:  500 * 1024 * 1024,
			}

			if err2 := d.GetDB().Create(&newBucket).Error; err2 != nil {
				return u.Message(false, "0x11134:Failed to create default bucket")
			}

			response := Response{
				BucketId:    newBucket.BucketId,
				BucketName:  newBucket.BucketName,
				EmptySize:   newBucket.EmptySize,
				SizeLimit:   newBucket.SizeLimit,
				BucketItems: []ResponceItem{},
			}
			resp := u.Message(true, "0x11006:Default bucket created")
			resp["buckets"] = response
			return resp
		}
	}

	var flatResults []flatResult

	sqlQuery := `
	WITH RECURSIVE BucketContent AS (
	    SELECT
	        bi.bucket_item_id,
	        bi.bucket_id,
	        bi.item_type,
	        bi.item_size,
	        bi.item_name,
	        bi.parent_id,
	        bi.item_url,
	        1 AS depth
	    FROM services.bucket_items bi
	    JOIN services.buckets b ON bi.bucket_id = b.bucket_id
	    WHERE b.user_id = ?
	      AND b.merchant_id = ?
	      AND b.app_id = ?
	      AND b.role_id = ?
	      AND (
	          ? IS NULL OR ? = '00000000-0000-0000-0000-000000000000' OR bi.bucket_id = ?
	      )
	      AND (
	          (? IS NULL OR ? = '00000000-0000-0000-0000-000000000000') AND bi.parent_id IS NULL
	          OR
	          (? IS NOT NULL AND ? != '00000000-0000-0000-0000-000000000000') AND bi.bucket_item_id = ?
	      )

	    UNION ALL

	    SELECT
	        bi.bucket_item_id,
	        bi.bucket_id,
	        bi.item_type,
	        bi.item_size,
	        bi.item_name,
	        bi.parent_id,
	        bi.item_url,
	        bc.depth + 1 AS depth
	    FROM services.bucket_items bi
	    INNER JOIN BucketContent bc ON bi.parent_id = bc.bucket_item_id
	    WHERE bc.depth < ?
	      AND (
	          ? IS NULL OR ? = '00000000-0000-0000-0000-000000000000' OR bi.bucket_id = ?
	      )
	)
	SELECT
	    bc.bucket_item_id,
	    bc.bucket_id,
	    b.bucket_name,
	    b.empty_size,
	    b.size_limit,
	    bc.item_name,
	    bc.item_type,
	    bc.item_size,
	    bc.item_url,
	    bc.parent_id,
	    bc.depth
	FROM BucketContent bc
	JOIN services.buckets b ON bc.bucket_id = b.bucket_id
	ORDER BY bc.bucket_id, bc.parent_id ASC NULLS FIRST, bc.item_type DESC, bc.item_name ASC;
	`

	err := d.GetDB().Raw(sqlQuery,
		userId,
		merchantId,
		appId,
		roleId,

		bucketId,
		bucketId,
		bucketId,

		bucketItemId,
		bucketItemId,

		bucketItemId,
		bucketItemId,
		bucketItemId,

		maxDepth,

		bucketId,
		bucketId,
		bucketId,
	).Scan(&flatResults).Error

	if err != nil {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	buckets := BuildNestedBuckets(flatResults)

	resp := u.Message(true, "0x11006:Requirement passed")
	resp["buckets"] = buckets

	return resp
}

// BuildNestedBuckets, flat sonuç listesini nested JSON formatına çevirir (recursive)
func BuildNestedBuckets(flatResults []flatResult) []Response {

	bucketsMap := make(map[uuid.UUID]*Response)
	itemsByBucket := make(map[uuid.UUID]map[uuid.UUID]*ResponceItem) // bucketId -> bucketItemId -> *ResponceItem

	for _, row := range flatResults {
		b, exists := bucketsMap[row.BucketId]
		if !exists {
			b = &Response{
				BucketId:    row.BucketId,
				BucketName:  row.BucketName,
				EmptySize:   row.EmptySize,
				SizeLimit:   row.SizeLimit,
				BucketItems: []ResponceItem{},
			}
			bucketsMap[row.BucketId] = b
			itemsByBucket[row.BucketId] = make(map[uuid.UUID]*ResponceItem)
		}

		parentId := uuid.Nil
		if row.ParentId != nil {
			parentId = *row.ParentId
		}

		item := &ResponceItem{
			BucketItemId: row.BucketItemId,
			ParentId:     parentId,
			ItemType:     row.ItemType,
			ItemName:     row.ItemName,
			ItemSize:     row.ItemSize,
			ItemURL:      row.ItemURL,
			ItemCount:    0,
			TotalSize:    0,
			BucketItems:  []ResponceItem{},
		}

		itemsByBucket[row.BucketId][row.BucketItemId] = item
	}

	for bucketId, itemsMap := range itemsByBucket {
		childrenMap := make(map[uuid.UUID][]*ResponceItem)
		var roots []*ResponceItem

		for _, item := range itemsMap {
			if item.ParentId == uuid.Nil {
				roots = append(roots, item)
			} else {
				childrenMap[item.ParentId] = append(childrenMap[item.ParentId], item)
			}
		}

		var buildTree func(node *ResponceItem) (int, float64)
		buildTree = func(node *ResponceItem) (int, float64) {
			children, found := childrenMap[node.BucketItemId]
			if !found {
				node.ItemCount = 1
				node.TotalSize = node.ItemSize
				return 1, node.ItemSize
			}

			count := 0
			size := 0.0
			for _, child := range children {
				cCount, cSize := buildTree(child)
				count += cCount
				size += cSize
				node.BucketItems = append(node.BucketItems, *child)
			}

			count += 1
			size += node.ItemSize

			node.ItemCount = count
			node.TotalSize = size

			return count, size
		}

		for _, root := range roots {
			buildTree(root)
			bucketsMap[bucketId].BucketItems = append(bucketsMap[bucketId].BucketItems, *root)
		}
	}

	buckets := make([]Response, 0, len(bucketsMap))
	for _, b := range bucketsMap {
		buckets = append(buckets, *b)
	}

	return buckets
}

func (req *RequestDto) Create() map[string]interface{} {

	// get bucket info by token info

	req.BucketItemId = uuid.New()

	bucketItem := &BucketItem{}

	bucketItem.BucketId = req.BucketId
	bucketItem.BucketItemId = req.BucketItemId
	bucketItem.ItemName = req.ItemName
	bucketItem.ItemSize = req.TotalSize
	bucketItem.ItemType = req.ItemType
	bucketItem.ItemURL = req.ItemURL

	d.GetDB().Table("application.notes").Create(bucketItem)

	if bucketItem.ID <= 0 {
		return u.Message(false, "0x11132:Failed to create note, connection error")
	}

	bucketItem.ID = 0

	response := u.Message(true, "0x11133:Note has been created")
	response["drive"] = req
	return response
}

func SoftDeleteBucketItemRecursive(
	bucketItemId uuid.UUID,
	userId uuid.UUID,
	merchantId uuid.UUID,
	appId int,
	roleId int,
) map[string]interface{} {

	// 1. Öncelikle, bucket_item_id ile başlayıp, ilgili bucket_id'yi ve user/app/merchant/role eşleşmesini kontrol et
	var bucketId uuid.UUID
	checkQuery := `
		SELECT b.bucket_id
		FROM services.bucket_items bi
		JOIN services.buckets b ON bi.bucket_id = b.bucket_id
		WHERE bi.bucket_item_id = ?
		  AND b.user_id = ?
		  AND b.merchant_id = ?
		  AND b.app_id = ?
		  AND b.role_id = ?
		LIMIT 1;
	`

	err := d.GetDB().Raw(checkQuery, bucketItemId, userId, merchantId, appId, roleId).Scan(&bucketId).Error
	if err != nil {
		return u.Message(false, fmt.Sprint("bucket item kontrolü sırasında hata: %w", err))
	}
	if bucketId == uuid.Nil {
		return u.Message(false, fmt.Sprint("bucket item bulunamadı veya erişim yetkiniz yok"))
	}

	// 2. Recursive SQL ile silinecek tüm bucket_item_id'leri bul
	var itemIDs []uuid.UUID
	recursiveQuery := `
	WITH RECURSIVE ItemsToDelete AS (
		SELECT bucket_item_id
		FROM services.bucket_items
		WHERE bucket_item_id = ?

		UNION ALL

		SELECT bi.bucket_item_id
		FROM services.bucket_items bi
		INNER JOIN ItemsToDelete itd ON bi.parent_id = itd.bucket_item_id
	)
	SELECT bucket_item_id FROM ItemsToDelete;
	`

	err = d.GetDB().Raw(recursiveQuery, bucketItemId).Scan(&itemIDs).Error
	if err != nil {
		return u.Message(false, fmt.Sprint("silinecek öğe id'leri alınırken hata: %w", err))
	}

	if len(itemIDs) == 0 {
		return u.Message(false, fmt.Sprint("silinecek öğe bulunamadı"))
	}

	// 3. Soft delete işlemini yap (gorm.Delete soft delete kullanır)
	// Transaction ile güvenli işlem
	err = d.GetDB().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("bucket_item_id IN ?", itemIDs).Delete(&BucketItem{})
		if result.Error != nil {
			return fmt.Errorf("soft delete sırasında hata: %w", result.Error)
		}
		return nil
	})

	if err != nil {
		return u.Message(false, fmt.Sprint("işlem hatası: %w", err))
	}

	return u.Message(true, "işlem basarili")
}

func HardDeleteBucketItemRecursive(
	bucketItemId uuid.UUID,
	userId uuid.UUID,
	merchantId uuid.UUID,
	appId int,
	roleId int,
) map[string]interface{} {

	// 1. Bucket ve kullanıcı bilgilerini kontrol et
	var bucketId uuid.UUID
	checkQuery := `
		SELECT b.bucket_id
		FROM services.bucket_items bi
		JOIN services.buckets b ON bi.bucket_id = b.bucket_id
		WHERE bi.bucket_item_id = ?
		  AND b.user_id = ?
		  AND b.merchant_id = ?
		  AND b.app_id = ?
		  AND b.role_id = ?
		LIMIT 1;
	`
	err := d.GetDB().Raw(checkQuery, bucketItemId, userId, merchantId, appId, roleId).Scan(&bucketId).Error
	if err != nil {
		return u.Message(false, fmt.Sprint("bucket item kontrolü sırasında hata: %w", err))
	}
	if bucketId == uuid.Nil {
		return u.Message(false, fmt.Sprint("bucket item bulunamadı veya erişim yetkiniz yok"))
	}

	// 2. Recursive alt öğeleri al
	var itemIDs []uuid.UUID
	recursiveQuery := `
	WITH RECURSIVE ItemsToDelete AS (
		SELECT bucket_item_id
		FROM services.bucket_items
		WHERE bucket_item_id = ?

		UNION ALL

		SELECT bi.bucket_item_id
		FROM services.bucket_items bi
		INNER JOIN ItemsToDelete itd ON bi.parent_id = itd.bucket_item_id
	)
	SELECT bucket_item_id FROM ItemsToDelete;
	`
	err = d.GetDB().Raw(recursiveQuery, bucketItemId).Scan(&itemIDs).Error
	if err != nil {
		return u.Message(false, fmt.Sprint("silinecek öğe id'leri alınırken hata: %w", err))
	}
	if len(itemIDs) == 0 {
		return u.Message(false, fmt.Sprint("silinecek öğe bulunamadı"))
	}

	// 3. Hard delete (kalıcı silme) işlemi
	err = d.GetDB().Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().Where("bucket_item_id IN ?", itemIDs).Delete(&BucketItem{})
		if result.Error != nil {
			return fmt.Errorf("hard delete sırasında hata: %w", result.Error)
		}
		return nil
	})

	if err != nil {
		return u.Message(false, fmt.Sprint("işlem hatası: %w", err))
	}

	return u.Message(true, "Islem basarili")
}

// func GetOrCreateBucketsWithItems(
// 	userId uuid.UUID,
// 	merchantId uuid.UUID,
// 	appId int,
// 	roleId int,
// ) ([]Response, map[string]interface{}, bool) {

// 	var bucketsCount int64
// 	err := d.GetDB().Table("services.buckets").
// 		Where("user_id = ? AND merchant_id = ? AND app_id = ? AND role_id = ?", userId, merchantId, appId, roleId).
// 		Count(&bucketsCount).Error
// 	if err != nil {
// 		return nil, u.Message(false, "0x11004:Connection error. Please retry"), false
// 	}

// 	// Bucket yoksa oluştur
// 	if bucketsCount == 0 {
// 		newBucket := Bucket{
// 			BucketId:   uuid.New(),
// 			UserId:     userId,
// 			RoleId:     roleId,
// 			AppId:      appId,
// 			MerchantId: merchantId,
// 			BucketName: "Default Bucket",
// 			EmptySize:  500 * 1024 * 1024, // 500 MB byte cinsinden
// 			SizeLimit:  500 * 1024 * 1024, // 500 MB byte cinsinden
// 		}

// 		err := d.GetDB().Create(&newBucket).Error
// 		if err != nil {
// 			return nil, u.Message(false, "0x11134:Failed to create default bucket"), false
// 		}

// 		response := Response{
// 			BucketId:    newBucket.BucketId,
// 			BucketName:  newBucket.BucketName,
// 			EmptySize:   newBucket.EmptySize,
// 			SizeLimit:   newBucket.SizeLimit,
// 			BucketItems: []ResponceItem{},
// 		}
// 		return []Response{response}, u.Message(true, "0x11006:Default bucket created"), true
// 	}
// }
