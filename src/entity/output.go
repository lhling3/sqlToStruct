package entity

import (
	"context"
	"gorm.io/gorm"
)

type CommercialUseRecord struct {
	Id          uint   `json:"id" gorm:"column:id"`
	AuthId      uint64 `json:"auth_id" gorm:"column:auth_id"`           //授权id
	SkuId       uint64 `json:"sku_id" gorm:"column:sku_id"`             //资源id
	PrivilegeId string `json:"privilege_id" gorm:"column:privilege_id"` //权益id
	SkuType     uint8  `json:"sku_type" gorm:"column:sku_type"`         //1:图片；2：图标；3：模板；4：字体
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`           //用户id
	Nickname    string `json:"nickname" gorm:"column:nickname"`         //名称
	Thumburl    string `json:"thumburl" gorm:"column:thumburl"`         //资源缩略图
	CreateTime  uint   `json:"create_time" gorm:"column:create_time"`   //创建时间
	UpdateTime  uint   `json:"update_time" gorm:"column:update_time"`   //修改时间
}

func (m *CommercialUseRecord) TableName() string {
	return "commercial_use_record"
}

func (m *CommercialUseRecord) Add(db *gorm.DB) (err error) {
	return helpers.WrapError(db.Table(m.TableName()).Create(m).Error)
}

func (m *CommercialUseRecord) Update(db *gorm.DB) (err error) {
	return helpers.WrapError(db.Table(m.TableName()).Updates(m).Error)
}

func (m *CommercialUseRecord) First(ctx context.Context) (err error) {
	return helpers.WrapError(global.DB(ctx, m.DbName()).Where(m).First(m).Error)
}

func (m *CommercialUseRecord) Lists(ctx context.Context, lists *[]CommercialUseRecord) (err error) {
	return helpers.WrapError(global.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Order("id desc").Find(lists).Error)
}

func (m *CommercialUseRecord) Count(ctx context.Context, count *int64) (err error) {
	return helpers.WrapError(global.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Count(count).Error)
}
