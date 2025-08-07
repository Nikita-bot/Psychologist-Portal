package handler

import "psyportal/dal"

type Handler struct {
	repo *dal.Repo
}

func InitHandler(repo *dal.Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}
