package service

import "server-go/internal/app/dao"

var GroupMemberService = &groupMemberService{}

type groupMemberService struct{}

// Count @Title 返回群聊人员总数
func (g *groupMemberService) Count(groupID int64) int {
	return int(dao.GroupMember.CountMember(groupID))
}
