package commands

import (
	"context"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/forms/domain"
)

type CreateGroupCommand struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

type CreateGroupHandler interface {
	Handle(ctx context.Context, cmd CreateGroupCommand) (*domain.FormGroup, error)
}

type createGroupHandler struct {
	repository domain.Repository
}

func NewCreateGroupHandler(repo domain.Repository) CreateGroupHandler {
	return &createGroupHandler{repository: repo}
}

func (h *createGroupHandler) Handle(ctx context.Context, cmd CreateGroupCommand) (*domain.FormGroup, error) {
	group := domain.NewFormGroup(cmd.Name, cmd.SortOrder)
	if err := h.repository.CreateGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

type UpdateGroupCommand struct {
	ID        string `json:"-"`
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

type UpdateGroupHandler interface {
	Handle(ctx context.Context, cmd UpdateGroupCommand) (*domain.FormGroup, error)
}

type updateGroupHandler struct {
	repository domain.Repository
}

func NewUpdateGroupHandler(repo domain.Repository) UpdateGroupHandler {
	return &updateGroupHandler{repository: repo}
}

func (h *updateGroupHandler) Handle(ctx context.Context, cmd UpdateGroupCommand) (*domain.FormGroup, error) {
	group, err := h.repository.GetGroupByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	group.Name = cmd.Name
	group.SortOrder = cmd.SortOrder
	group.UpdatedAt = time.Now()
	if err := h.repository.UpdateGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

type DeleteGroupCommand struct {
	ID string
}

type DeleteGroupHandler interface {
	Handle(ctx context.Context, cmd DeleteGroupCommand) error
}

type deleteGroupHandler struct {
	repository domain.Repository
}

func NewDeleteGroupHandler(repo domain.Repository) DeleteGroupHandler {
	return &deleteGroupHandler{repository: repo}
}

func (h *deleteGroupHandler) Handle(ctx context.Context, cmd DeleteGroupCommand) error {
	return h.repository.DeleteGroup(ctx, cmd.ID)
}
