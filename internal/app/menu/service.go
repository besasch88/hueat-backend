package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"gorm.io/gorm"
)

type menuServiceInterface interface {
	getMenu(ctx *gin.Context, input getMenuInputDto) (menu, error)
}

type menuService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  menuRepositoryInterface
}

func newMenuService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository menuRepositoryInterface) menuService {
	return menuService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuService) getMenu(ctx *gin.Context, input getMenuInputDto) (menu, error) {
	targetTableID := hueat_utils.GetUUIDFromString(input.TableID)
	table, err := s.repository.getTableByID(s.storage, targetTableID)
	if err != nil {
		return menu{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(table) {
		return menu{}, errTableNotFound
	}
	categories, _, err := s.repository.listMenuCategories(s.storage, *table.Inside, false)
	if err != nil || categories == nil {
		return menu{}, hueat_err.ErrGeneric
	}
	items, _, err := s.repository.listMenuItems(s.storage, *table.Inside, targetTableID, false)
	if err != nil || items == nil {
		return menu{}, hueat_err.ErrGeneric
	}
	options, _, err := s.repository.listMenuOptions(s.storage, *table.Inside, false)
	if err != nil || options == nil {
		return menu{}, hueat_err.ErrGeneric
	}

	menu := menu{
		Categories: []menuCategory{},
	}
	for _, category := range categories {
		menuCategory := menuCategory{
			menuCategoryEntity: category,
			Items:              []menuItem{},
		}
		for _, item := range items {
			menuItem := menuItem{
				menuItemEntity: item,
				Options:        []menuOption{},
			}
			if item.MenuCategoryID.String() != category.ID.String() {
				continue
			}
			for _, option := range options {
				if option.MenuItemID.String() != item.ID.String() {
					continue
				}
				menuOption := menuOption{
					menuOptionEntity: option,
				}
				menuItem.Options = append(menuItem.Options, menuOption)
			}
			menuCategory.Items = append(menuCategory.Items, menuItem)
		}
		menu.Categories = append(menu.Categories, menuCategory)
	}
	return menu, nil
}
