package blocks

import (
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
)

type BlockService struct {
	DB *pg.DB
}

// Get block by height (id)
func (service *BlockService) GetById(id uint64) *BlockResource {
	// fetch model from the database
	var block models.Block
	err := service.DB.Model(&block).Column("Validators").Where("ID = ?", id).Select()

	// check to existing row
	if err != nil {
		return nil
	}

	transformedBlock := TransformBlock(block)
	return &transformedBlock
}

// Get paginated list of blocks
func (service *BlockService) GetList(page int, perPage int) *[]BlockResource {
	var blocks []models.Block

	// fetch blocks
	err := service.DB.Model(&blocks).Column("Validators").Limit(perPage).
		Offset(perPage * (page - 1)).Order("id DESC").Select()

	// check to existing
	if (err != nil) || (len(blocks) == 0) {
		empty := make([]BlockResource, 0)
		return &empty
	}

	// transform models
	var blocksResource []BlockResource
	for _, block := range blocks {
		blocksResource = append(blocksResource, TransformBlock(block))
	}

	return &blocksResource
}
