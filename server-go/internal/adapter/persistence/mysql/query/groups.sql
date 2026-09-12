-- name: InsertGroup :execresult
INSERT INTO `groups` (name, avatar, owner_id) VALUES (?, ?, ?);

-- name: GetGroupByID :one
SELECT * FROM `groups` WHERE id = ?;

-- name: InsertGroupMember :exec
INSERT IGNORE INTO group_members (group_id, user_id, role) VALUES (?, ?, ?);

-- name: CountMembership :one
SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?;

-- name: ListGroupMemberIDs :many
SELECT user_id FROM group_members WHERE group_id = ?;

-- name: CountGroupMembers :one
SELECT COUNT(*) FROM group_members WHERE group_id = ?;

-- name: ListGroupsByUser :many
SELECT g.* FROM `groups` g
JOIN group_members m ON g.id = m.group_id
WHERE m.user_id = ?;

-- name: SearchGroupsByName :many
SELECT * FROM `groups` WHERE name LIKE CONCAT(?, '%') LIMIT ?;
