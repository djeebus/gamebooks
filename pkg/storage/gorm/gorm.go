package gorm

import (
	"encoding/json"
	"fmt"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func New(dialect gorm.Dialector) (*Gorm, error) {
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	if err = db.AutoMigrate(new(UserStorage)); err != nil {
		return nil, errors.Wrap(err, "failed to migrate UserStorage")
	}

	return &Gorm{db: db}, nil
}

type Gorm struct {
	db *gorm.DB
}

var _ storage.Storage = new(Gorm)

func (g *Gorm) marshal(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	return encoded, errors.Wrap(err, "failed to encode")
}

func (g *Gorm) unmarshal(data []byte, value any) error {
	err := json.Unmarshal(data, value)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	return nil
}

func (g *Gorm) Get(key string) (interface{}, error) {
	var model UserStorage
	if err := g.db.Find(&model, "key = ?", key).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find key")
	}

	if model.Value == nil {
		return nil, nil
	}

	var data any
	if err := g.unmarshal(model.Value, &data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal key")
	}

	return data, nil
}

func (g *Gorm) Set(key string, value interface{}) error {
	encoded, err := g.marshal(value)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	row := UserStorage{
		Key:   key,
		Value: encoded,
	}

	if err = g.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"key", "value"}),
	}).Create(&row).Error; err != nil {
		return errors.Wrap(err, "failed to set")
	}

	return nil
}

func (g *Gorm) Remove(key string) error {
	model := UserStorage{
		Key: key,
	}

	if err := g.db.Delete(&model).Error; err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	return nil
}

func (g *Gorm) Clear(keyPrefix string) error {
	if err := g.db.Delete(
		new(UserStorage),
		"key LIKE ?", fmt.Sprintf("%s%%", keyPrefix),
	).Error; err != nil {
		return errors.Wrap(err, "failed to delete keys")
	}

	return nil
}
