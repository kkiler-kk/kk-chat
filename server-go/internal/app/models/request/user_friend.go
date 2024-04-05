package request

type UserFriendAddReqInput struct {
	UserId   int64 `from:"userId" binding:"required"`
	FriendId int64 `from:"friendId" binding:"required"`
}
