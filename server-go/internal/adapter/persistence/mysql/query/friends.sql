-- name: InsertFriend :exec
INSERT IGNORE INTO friends (user_id, friend_id) VALUES (?, ?);

-- name: CountFriendship :one
SELECT COUNT(*) FROM friends WHERE user_id = ? AND friend_id = ?;

-- name: ListFriendIDs :many
SELECT friend_id FROM friends WHERE user_id = ?;
