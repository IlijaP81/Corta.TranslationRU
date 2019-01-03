package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	MessageRepository interface {
		With(ctx context.Context, db *factory.DB) MessageRepository

		FindMessageByID(id uint64) (*types.Message, error)
		FindMessages(filter *types.MessageFilter) (types.MessageSet, error)
		FindThreads(filter *types.MessageFilter) (types.MessageSet, error)
		CountFromMessageID(channelID, threadID, messageID uint64) (uint32, error)
		PrefillThreadParticipants(mm types.MessageSet) error
		CreateMessage(mod *types.Message) (*types.Message, error)
		UpdateMessage(mod *types.Message) (*types.Message, error)
		DeleteMessageByID(ID uint64) error
		IncReplyCount(ID uint64) error
		DecReplyCount(ID uint64) error
	}

	message struct {
		*repository
	}
)

const (
	MESSAGES_MAX_LIMIT = 100

	sqlMessageColumns = "id, " +
		"COALESCE(type,'') AS type, " +
		"message, " +
		"rel_user, " +
		"rel_channel, " +
		"reply_to, " +
		"replies, " +
		"created_at, " +
		"updated_at, " +
		"deleted_at"
	sqlMessageScope = "deleted_at IS NULL"

	sqlMessagesSelect = `SELECT ` + sqlMessageColumns + `
        FROM messages
       WHERE ` + sqlMessageScope

	sqlMessagesThreads = "WITH originals AS (" +
		" SELECT id AS original_id " +
		"   FROM messages " +
		"  WHERE " + sqlMessageScope +
		"    AND rel_channel IN " + sqlChannelAccess +
		"    AND reply_to = 0 " +
		"    AND replies > 0 " +
		// for finding only threads we've created or replied to
		"    AND (rel_user = ? OR id IN (SELECT DISTINCT reply_to FROM messages WHERE rel_user = ?))" +
		"  ORDER BY id DESC " +
		"  LIMIT ? " +
		")" +
		" SELECT " + sqlMessageColumns +
		"   FROM messages, originals " +
		"  WHERE " + sqlMessageScope +
		"    AND original_id IN (id, reply_to)"

	sqlThreadParticipantsByMessageID = "SELECT DISTINCT reply_to, rel_user FROM messages WHERE reply_to IN (?)"

	sqlCountFromMessageID = "SELECT COUNT(*) AS count FROM messages WHERE rel_channel = ? AND reply_to = ? AND id > ?"

	sqlMessageRepliesIncCount = `UPDATE messages SET replies = replies + 1 WHERE id = ? AND reply_to = 0`
	sqlMessageRepliesDecCount = `UPDATE messages SET replies = replies - 1 WHERE id = ? AND reply_to = 0`

	ErrMessageNotFound = repositoryError("MessageNotFound")
)

func Message(ctx context.Context, db *factory.DB) MessageRepository {
	return (&message{}).With(ctx, db)
}

func (r *message) With(ctx context.Context, db *factory.DB) MessageRepository {
	return &message{
		repository: r.repository.With(ctx, db),
	}
}

func (r *message) FindMessageByID(id uint64) (*types.Message, error) {
	mod := &types.Message{}
	sql := sqlMessagesSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrMessageNotFound)
}

func (r *message) FindMessages(filter *types.MessageFilter) (types.MessageSet, error) {
	r.sanitizeFilter(filter)

	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	sql := sqlMessagesSelect

	if filter.Query != "" {
		sql += " AND message LIKE ?"
		params = append(params, "%"+filter.Query+"%")
	}

	if filter.ChannelID > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelID)
	}

	if filter.UserID > 0 {
		sql += " AND rel_user = ? "
		params = append(params, filter.UserID)
	}

	if filter.RepliesTo > 0 {
		sql += " AND reply_to = ? "
		params = append(params, filter.RepliesTo)
	} else {
		sql += " AND reply_to = 0 "
	}

	// first, exclusive
	if filter.FirstID > 0 {
		sql += " AND id > ? "
		params = append(params, filter.FirstID)
	}

	// from, inclusive
	if filter.FromID > 0 {
		sql += " AND id >= ? "
		params = append(params, filter.FromID)
	}

	// last, exclusive
	if filter.LastID > 0 {
		sql += " AND id < ? "
		params = append(params, filter.LastID)
	}

	// to, inclusive
	if filter.ToID > 0 {
		sql += " AND id <= ? "
		params = append(params, filter.ToID)
	}

	if filter.BookmarkedOnly || filter.PinnedOnly {
		sql += " AND id IN (SELECT rel_message FROM message_flags WHERE flag = ?) "

		if filter.PinnedOnly {
			params = append(params, types.MessageFlagBookmarkedMessage)
		} else {
			params = append(params, types.MessageFlagPinnedToChannel)
		}
	}

	sql += " AND rel_channel IN " + sqlChannelAccess
	params = append(params, filter.CurrentUserID, types.ChannelTypePublic)

	sql += " ORDER BY id DESC"

	sql += " LIMIT ? "
	params = append(params, filter.Limit)

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *message) FindThreads(filter *types.MessageFilter) (types.MessageSet, error) {
	r.sanitizeFilter(filter)

	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	// for sqlChannelAccess
	params = append(params, filter.CurrentUserID, types.ChannelTypePublic)

	// for finding only threads we've created or replied to
	params = append(params, filter.CurrentUserID, filter.CurrentUserID)

	// for sqlMessagesThreads
	params = append(params, filter.Limit)

	sql := sqlMessagesThreads
	if filter.ChannelID > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelID)
	}

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *message) CountFromMessageID(channelID, threadID, messageID uint64) (uint32, error) {
	rval := struct{ Count uint32 }{}
	return rval.Count, r.db().Get(&rval, sqlCountFromMessageID, channelID, threadID, messageID)
}

func (r *message) PrefillThreadParticipants(mm types.MessageSet) error {
	var rval = []struct {
		ReplyTo uint64 `db:"reply_to"`
		UserID  uint64 `db:"rel_user"`
	}{}

	if len(mm) == 0 {
		return nil
	}

	if sql, args, err := sqlx.In(sqlThreadParticipantsByMessageID, mm.IDs()); err != nil {
		return err
	} else if err = r.db().Select(&rval, sql, args...); err != nil {
		return err
	} else {
		for _, p := range rval {
			mm.FindByID(p.ReplyTo).RepliesFrom = append(mm.FindByID(p.ReplyTo).RepliesFrom, p.UserID)
		}
	}

	return nil
}

func (r *message) sanitizeFilter(filter *types.MessageFilter) {
	if filter == nil {
		filter = &types.MessageFilter{}
	}

	if filter.Limit == 0 || filter.Limit > MESSAGES_MAX_LIMIT {
		filter.Limit = MESSAGES_MAX_LIMIT
	}
}

func (r *message) CreateMessage(mod *types.Message) (*types.Message, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("messages", mod)
}

func (r *message) UpdateMessage(mod *types.Message) (*types.Message, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace("messages", mod)
}

func (r *message) DeleteMessageByID(ID uint64) error {
	return r.updateColumnByID("messages", "deleted_at", time.Now(), ID)
}

func (r *message) IncReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesIncCount, ID)
	return err
}

func (r *message) DecReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesDecCount, ID)
	return err
}
