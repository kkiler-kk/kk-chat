package mongox

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// 集合 "messages" 文档结构：
//
//	{ _id, conversation_id: "u:1_2", type: "private", sender_id: NumberLong,
//	  content: "...", content_type: "text", created_at: ISODate }
type messageDoc struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ConversationID string             `bson:"conversation_id"`
	Type           string             `bson:"type"`
	SenderID       int64              `bson:"sender_id"`
	Content        string             `bson:"content"`
	ContentType    string             `bson:"content_type"`
	CreatedAt      time.Time          `bson:"created_at"`
}

type MessageRepo struct {
	col           *mongo.Collection
	retentionDays int
}

func NewMessageRepo(db *mongo.Database, retentionDays int) *MessageRepo {
	return &MessageRepo{col: db.Collection("messages"), retentionDays: retentionDays}
}

// EnsureIndexes 幂等创建：查询索引 (conversation_id, created_at DESC) + TTL 索引。
func (r *MessageRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "conversation_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(int32(r.retentionDays * 86400)),
		},
	})
	return err
}

func (r *MessageRepo) Save(ctx context.Context, msg *domain.Message) error {
	doc := messageDoc{
		ConversationID: msg.ConversationID.Value,
		Type:           msg.ConversationID.Type.String(),
		SenderID:       msg.SenderID,
		Content:        msg.Content,
		ContentType:    string(msg.ContentType),
		CreatedAt:      msg.CreatedAt.UTC(),
	}
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		msg.ID = oid.Hex()
	}
	return nil
}

// History 取 created_at < cursor 的最近 limit 条（倒序查询后反转为升序返回）；
// cursor 零值表示不加时间过滤。
func (r *MessageRepo) History(ctx context.Context, convID domain.ConversationID, cursor time.Time, limit int) ([]domain.Message, error) {
	filter := bson.M{"conversation_id": convID.Value}
	if !cursor.IsZero() {
		filter["created_at"] = bson.M{"$lt": cursor.UTC()}
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var docs []messageDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	out := make([]domain.Message, 0, len(docs))
	for i := len(docs) - 1; i >= 0; i-- { // 反转为升序
		d := docs[i]
		out = append(out, domain.Message{
			ID:             d.ID.Hex(),
			ConversationID: convID,
			SenderID:       d.SenderID,
			Content:        d.Content,
			ContentType:    domain.ContentType(d.ContentType),
			CreatedAt:      d.CreatedAt,
		})
	}
	return out, nil
}

var _ port.MessageRepo = (*MessageRepo)(nil)
