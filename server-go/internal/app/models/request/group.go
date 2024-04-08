package request

type GroupCreateReqInput struct {
	Name      string  `from:"name" binding:"required"`
	OwnerId   int64   `from:"ownerId" binding:"required"`
	Icon      string  `from:"icon"`
	MembersID []int64 `from:"membersId" binding:"required"`
}
