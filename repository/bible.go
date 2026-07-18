package repository

import "github.com/ian77-huang/golang-echo/model"

func (r *bibleRepository) GetBible(id int) (*model.Bible, error) {
	bible := &model.Bible{}
	if err := r.db.Where("id = ?", id).First(bible).Error; err != nil {
		return nil, err
	}
	return bible, nil
}
