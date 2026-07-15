package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/community"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/server"
)

type CommunityRepository struct {
	server *server.Server
}

func NewCommunityRepository(server *server.Server) *CommunityRepository {
	return &CommunityRepository{
		server: server,
	}
}

func (r *CommunityRepository) CreateCommunity(ctx context.Context, payload *community.CreateCommunityPayload) (*community.Community, error) {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	stmt := `
		INSERT INTO communities(
			admin_id,
			name,
			slug,
			description,
			avatar_key,
			banner_key
		)
		VALUES(
			@admin_id,
			@name,
			@slug,
			@description,
			@avatar_key,
			@banner_key
		)
		RETURNING *
	`
	row, err := tx.Query(ctx, stmt, pgx.NamedArgs{
		"admin_id":    payload.AdminID,
		"name":        payload.Name,
		"slug":        payload.Slug,
		"description": payload.Description,
		"avatar_key":  payload.AvatarKey,
		"banner_key":  payload.BannerKey,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create community: %w", err)
	}
	newCommunity, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[community.Community])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to struct: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO community_members(user_id, community_id, role)
		VALUES(@user_id, @community_id, @role)
		`,
		pgx.NamedArgs{
			"user_id":      payload.AdminID,
			"community_id": newCommunity.ID,
			"role":         community.AdminRole,
		})
	if err != nil {
		return nil, fmt.Errorf("Failed to add admin as community member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &newCommunity, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) UpdateCommunitySettings(ctx context.Context, communityID uuid.UUID, payload *community.UpdateCommunitySettingsPayload) (*community.Community, error) {
	check, err := r.IsModerator(ctx, communityID, payload.UserID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", payload.UserID, communityID, err)
	}

	if *check == false {
		return nil, fmt.Errorf("The user with user_id %s is not a moderator/admin of the community", payload.UserID)
	}

	stmt := `
		UPDATE communities SET
	`
	setClauses := []string{}

	args := pgx.NamedArgs{
		"community_id": communityID,
	}

	if payload.Name != nil {
		args["name"] = *payload.Name
		setClauses = append(setClauses, "name = @name")
	}

	if payload.Slug != nil {
		args["slug"] = *payload.Slug
		setClauses = append(setClauses, "slug = @slug")
	}

	if payload.Description != nil {
		args["description"] = *payload.Description
		setClauses = append(setClauses, "description = @description")
	}

	if payload.AvatarKey != nil {
		args["avatar_key"] = *payload.AvatarKey
		setClauses = append(setClauses, "avatar_key = @avatar_key")
	}

	if payload.BannerKey != nil {
		args["banner_key"] = *payload.BannerKey
		setClauses = append(setClauses, "banner_key = @banner_key")
	}

	if len(setClauses) == 0 {
		code := "NOTHING_TO_UPDATE"
		return nil, errs.NewBadRequestError("No fields provided to update", false, &code, nil, nil)
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to update the community settings for community_id %s: %w", communityID, err)
	}

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.Community])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse updated community for community_id %s: %w", communityID, err)
	}

	return &updated, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) DeleteCommunity(ctx context.Context, communityID uuid.UUID, userID uuid.UUID) error {
	stmt := `
		DELETE FROM communities 
		WHERE id = @community_id
		AND admin_id = @user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"community_id": communityID,
		"user_id":      userID,
	})

	if err != nil {
		return fmt.Errorf("Failed to delete community with community_id %s: %w", communityID, err)
	}

	if result.RowsAffected() == 0 {
		code := "COMMUNITY_NOT_FOUND"
		return errs.NewNotFoundError("community not found", false, &code)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) ChangeMemberRoleInCommunity(ctx context.Context, communityID uuid.UUID, userID uuid.UUID, payload *community.ChangeMemberRoleInCommunityPayload) (*community.CommunityMember, error) {

	check, err := r.IsAdmin(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was the admin for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if !*check {
		code := "NOT_ADMIN"
		return nil, errs.NewBadRequestError("The user is not the admin of the community", false, &code, nil, nil)
	}

	if payload.TargetUserID == userID {
		code := "CANNOT_CHANGE_OWN_ROLE"
		return nil, errs.NewBadRequestError("cannot change your own role", false, &code, nil, nil)
	}

	stmt := `
		UPDATE community_members
		SET role = @role
		WHERE community_id = @community_id AND user_id = @target_user_id
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"role":           payload.NewRole,
		"community_id":   communityID,
		"target_user_id": payload.TargetUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to update the role of the member: %w", err)
	}

	communityMember, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.CommunityMember])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse updated member: %w", err)
	}

	return &communityMember, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetCommunityPostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, communityID uuid.UUID) (*post.PopulatedPost, error) {

	isUserBanned, err := r.IsBannedFromCommunity(ctx, userID, communityID)
	if err != nil {
		return nil, fmt.Errorf("Could Not check if the user was banned by the community: %w", err)
	}
	if *isUserBanned {
		code := "USER_IS_BANNED"
		return nil, errs.NewBadRequestError("user is banned", false, &code, nil, nil)
	}

	stmt := `
		SELECT p.*,
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY pm.created_at DESC, pm.id DESC
			) FILTER(WHERE pm.id IS NOT NULL),
			'[]'::JSONB
		) AS post_media,
		json_build_object(
			'communityId', c.id,
			'communityName', c.name,
			'communityAvatarKey', c.avatar_key
		) AS mini_community
		FROM posts p
		LEFT JOIN post_media pm ON pm.post_id = p.id
		LEFT JOIN communities c ON c.id = p.community_id
		WHERE p.id = @post_id
			AND p.community_id = @community_id
			AND p.deleted_at IS NULL
			AND NOT EXISTS (
				SELECT 1 FROM user_blocks ub
				WHERE (ub.blocker_id = @user_id AND ub.blocked_id = p.author_id)
				   OR (ub.blocker_id = p.author_id AND ub.blocked_id = @user_id)
			)
		GROUP BY p.id, c.id
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"post_id":      postID,
		"community_id": communityID,
		"user_id":      userID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to process get post by id query for post_id %s: %w", postID, err)
	}

	populatedPost, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[post.PopulatedPost])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to the struct: %w", err)
	}

	return &populatedPost, nil
}

// -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (r *CommunityRepository) GetCommunityMembers(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, query *community.GetCommunityMembersQuery) (*model.CursorPaginatedResponse[community.MiniCommunityUser], error) {

	isUserBanned, err := r.IsBannedFromCommunity(ctx, userID, communityID)
	if err != nil {
		return nil, fmt.Errorf("Could Not check if the user was banned by the community: %w", err)
	}
	if *isUserBanned {
		code := "USER_IS_BANNED"
		return nil, errs.NewBadRequestError("user is banned", false, &code, nil, nil)
	}

	stmt := `
		SELECT cm.user_id, u.avatar_key, u.display_name AS name, cm.joined_at, cm.role
		FROM community_members cm
		JOIN users u ON u.id = cm.user_id
		WHERE cm.community_id = @community_id
	`

	limit := 20
	args := pgx.NamedArgs{
		"community_id":   communityID,
		"limit_plus_one": limit + 1,
	}

	if query.Role != nil && *query.Role != community.CommonRole {
		stmt += " AND cm.role = @role"
		args["role"] = *query.Role
	}

	orderStmt := ""
	if *query.Order == model.OrderDesc {
		orderStmt = "ORDER BY cm.joined_at DESC"
		if query.CursorSortValue != nil {
			cursorTime, err := time.Parse(time.RFC3339Nano, *query.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cursor sort value: %w", err)
			}
			stmt += " AND cm.joined_at <= @joined_at"
			args["joined_at"] = cursorTime
		}
	} else {
		orderStmt = "ORDER BY cm.joined_at ASC"
		if query.CursorSortValue != nil {
			cursorTime, err := time.Parse(time.RFC3339Nano, *query.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cursor sort value: %w", err)
			}
			stmt += " AND cm.joined_at >= @joined_at"
			args["joined_at"] = cursorTime
		}
	}

	stmt += " " + orderStmt + " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to process get community members of community with comunity_id %s: %w", communityID, err)
	}

	communityMembers, err := pgx.CollectRows(rows, pgx.RowToStructByName[community.MiniCommunityUser])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the rows to the struct: %w", err)
	}

	if len(communityMembers) < limit+1 {
		var cursorSortValue string
		var cursorCreatedAt time.Time
		if query.CursorSortValue != nil {
			cursorSortValue = *query.CursorSortValue
		}
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[community.MiniCommunityUser]{
			Data:            communityMembers,
			CursorSortValue: cursorSortValue,
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	nextCursorSortValue := communityMembers[limit].JoinedAt.Format(time.RFC3339Nano)
	nextCursorCreatedAt := communityMembers[limit].JoinedAt

	return &model.CursorPaginatedResponse[community.MiniCommunityUser]{
		Data:            communityMembers[:limit],
		CursorSortValue: nextCursorSortValue,
		CursorCreatedAt: nextCursorCreatedAt,
		HasMore:         true,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) ReportCommunityPost(ctx context.Context, userID uuid.UUID, payload *community.ReportCommunityPostPayload) (*community.CommunityReport, error) {
	stmt := `
		INSERT INTO community_reports(reporter_id, community_id, post_id, reason, status)
		VALUES(@reporter_id, @community_id, @post_id, @reason, @status)
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"reporter_id":  userID,
		"community_id": payload.CommunityID,
		"post_id":      payload.PostID,
		"reason":       payload.Reason,
		"status":       community.ReportPending,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create a report: %w", err)
	}
	communityReport, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.CommunityReport])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to struct: %w", err)
	}
	return &communityReport, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) ResolveCommunityPostReport(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.ResolveCommunityPostReportPayload) (*community.CommunityReport, error) {
	check, err := r.IsModerator(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was the moderator/admin for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if !*check {
		code := "NOT_MODERATOR"
		return nil, errs.NewBadRequestError("user is not a moderator/admin of the community", false, &code, nil, nil)
	}

	stmt := `
		UPDATE community_reports
		SET status = @updated_status
		WHERE id = @report_id AND community_id = @community_id
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"report_id":      payload.ReportID,
		"updated_status": payload.UpdatedStatus,
		"community_id":   communityID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to update the report with report_id %s: %w", payload.ReportID, err)
	}
	report, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.CommunityReport])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse updated report: %w", err)
	}
	return &report, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) DeleteCommunityPost(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.DeleteCommunityPostPayload) error {
	check, err := r.IsModerator(ctx, communityID, userID)

	if err != nil {
		return fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if !*check {
		code := "NOT_MODERATOR"
		return errs.NewBadRequestError("user is not a moderator/admin of the community", false, &code, nil, nil)
	}

	stmt := `
		DELETE FROM posts
		WHERE id = @post_id
		AND community_id = @community_id
	`
	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"post_id":      payload.PostID,
		"community_id": communityID,
	})

	if err != nil {
		return fmt.Errorf("Failed to delete the post with post_id %s: %w", payload.PostID, err)
	}

	if result.RowsAffected() == 0 {
		code := "POST_NOT_FOUND"
		return errs.NewNotFoundError("post not found or does not belong to this community", false, &code)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) BanUserFromCommunity(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.BanCommunityMemberPayload) (*community.BannedFromCommunityUser, error) {
	check, err := r.IsModerator(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if !*check {
		code := "NOT_MODERATOR"
		return nil, errs.NewBadRequestError("user is not a moderator/admin of the community", false, &code, nil, nil)
	}

	if payload.UserIDToBan == userID {
		code := "CANNOT_BAN_SELF"
		return nil, errs.NewBadRequestError("cannot ban yourself", false, &code, nil, nil)
	}

	targetIsAdmin, err := r.IsAdmin(ctx, communityID, payload.UserIDToBan)
	if err != nil {
		return nil, fmt.Errorf("failed to check target admin status: %w", err)
	}
	if *targetIsAdmin {
		code := "CANNOT_BAN_ADMIN"
		return nil, errs.NewBadRequestError("cannot ban the community admin", false, &code, nil, nil)
	}

	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		INSERT INTO banned_from_community_users (community_id, user_id)
		VALUES (@community_id, @user_id_to_ban)
		RETURNING *
	`, pgx.NamedArgs{
		"community_id":   communityID,
		"user_id_to_ban": payload.UserIDToBan,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ban user %s from community %s by user %s: %w", payload.UserIDToBan, communityID, userID, err)
	}
	ban, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.BannedFromCommunityUser])
	if err != nil {
		return nil, fmt.Errorf("failed to parse banned user for community %s: %w", communityID, err)
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM community_members WHERE community_id = @community_id AND user_id = @user_id_to_ban
	`, pgx.NamedArgs{"community_id": communityID, "user_id_to_ban": payload.UserIDToBan})
	if err != nil {
		return nil, fmt.Errorf("failed to remove banned user from community_members: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &ban, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetCommunityReports(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, query *community.GetCommunityReportsQuery) (*model.CursorPaginatedResponse[community.CommunityReport], error) {

	check, err := r.IsModerator(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if !*check {
		code := "NOT_MODERATOR"
		return nil, errs.NewBadRequestError("user is not a moderator/admin of the community", false, &code, nil, nil)
	}

	stmt := `SELECT * FROM community_reports WHERE community_id = @community_id`
	limit := 20
	args := pgx.NamedArgs{
		"community_id":   communityID,
		"limit_plus_one": limit + 1,
	}

	conditions := []string{}

	if query.Status != nil {
		args["status"] = *query.Status
		conditions = append(conditions, "status = @status")
	}

	if query.ReportedDateStart != "" {
		start, err := time.Parse(time.RFC3339Nano, query.ReportedDateStart)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reported date start: %w", err)
		}
		args["reported_date_start"] = start
		conditions = append(conditions, "created_at >= @reported_date_start")
	}

	if query.ReportedDateEnd != "" {
		end, err := time.Parse(time.RFC3339Nano, query.ReportedDateEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reported date end: %w", err)
		}
		args["reported_date_end"] = end
		conditions = append(conditions, "created_at <= @reported_date_end")
	}

	if query.CursorCreatedAt != nil {
		args["cursor_created_at"] = *query.CursorCreatedAt
		conditions = append(conditions, "created_at <= @cursor_created_at")
	}

	if len(conditions) > 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	stmt += " ORDER BY created_at DESC LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query community reports for community_id %s: %w", communityID, err)
	}

	reports, err := pgx.CollectRows(rows, pgx.RowToStructByName[community.CommunityReport])
	if err != nil {
		return nil, fmt.Errorf("failed to collect community reports for community_id %s: %w", communityID, err)
	}

	if len(reports) < limit+1 {
		var cursorCreatedAt time.Time
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[community.CommunityReport]{
			Data:            reports,
			CursorSortValue: "",
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	return &model.CursorPaginatedResponse[community.CommunityReport]{
		Data:            reports[:limit],
		CursorSortValue: reports[limit].CreatedAt.Format(time.RFC3339Nano),
		CursorCreatedAt: reports[limit].CreatedAt,
		HasMore:         true,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetReportByID(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, reportID uuid.UUID) (*community.CommunityReport, error) {

	check, err := r.IsModerator(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if !*check {
		code := "NOT_MODERATOR"
		return nil, errs.NewBadRequestError("user is not a moderator/admin of the community", false, &code, nil, nil)
	}

	stmt := `
		SELECT * FROM community_reports
		WHERE id = @report_id AND community_id = @community_id
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"report_id":    reportID,
		"community_id": communityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query report_id %s in community_id %s: %w", reportID, communityID, err)
	}
	report, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.CommunityReport])
	if err != nil {
		return nil, fmt.Errorf("failed to parse report_id %s: %w", reportID, err)
	}
	return &report, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) IsBannedFromCommunity(ctx context.Context, userID uuid.UUID, communityID uuid.UUID) (*bool, error) {

	stmt := `
		SELECT EXISTS (
			SELECT 1
			FROM banned_from_community_users
			WHERE community_id = @community_id
			AND user_id = @user_id
		)
	`

	var banned bool

	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"community_id": communityID,
		"user_id":      userID,
	}).Scan(&banned)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check if user %s is banned from community %s: %w",
			userID,
			communityID,
			err,
		)
	}

	return &banned, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) IsModerator(ctx context.Context, communityID uuid.UUID, userID uuid.UUID) (*bool, error) {
	stmt := `
		SELECT EXISTS(
			SELECT 1
			FROM community_members cm
			WHERE cm.id = @community_id
			AND cm.user_id = @user_id
			AND role IN ('admin', 'moderator')
		)
	`

	var check bool

	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
	}).Scan(&check)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	return &check, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) IsAdmin(ctx context.Context, communityID uuid.UUID, userID uuid.UUID) (*bool, error) {
	stmt := `
		SELECT EXISTS(
			SELECT 1
			FROM community_members cm
			WHERE cm.id = @community_id
			AND cm.user_id = @user_id
			AND role = 'admin'
		)
	`

	var check bool

	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
	}).Scan(&check)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a admin for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	return &check, nil
}

func (r *CommunityRepository) IsMember(ctx context.Context, communityID uuid.UUID, userID uuid.UUID) (*bool, error) {
	stmt := `
		SELECT EXISTS(
			SELECT 1 FROM community_members WHERE community_id = @community_id AND user_id = @user_id
		)
	`

	var check bool

	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
	}).Scan(&check)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a member for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	return &check, nil
}

func (r *CommunityRepository) GetUserRole(ctx context.Context, communityID uuid.UUID, userID uuid.UUID) (*community.CommunityRole, error) {
	stmt := `
		SELECT role FROM community_members
		WHERE community_id = @community_id AND user_id = @user_id
	`
	var role community.CommunityRole
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
	}).Scan(&role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to check role for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	return &role, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetCommunities(ctx context.Context, query *community.GetCommunitiesQuery) (*model.CursorPaginatedResponse[community.MiniCommunity], error) {

	stmt := `
		SELECT c.id AS community_id, c.name, c.avatar_key
		FROM communities c
	`
	conditions := []string{}
	args := pgx.NamedArgs{}

	if query.Name != nil && *query.Name != "" {
		args["name"] = *query.Name + "%"
		conditions = append(conditions, "c.name ILIKE @name")
	}

	orderStmt := ""

	if *query.Sort == model.SortByMembersCount {
		if *query.Order == model.OrderDesc {
			orderStmt = "ORDER BY c.members_count DESC, c.created_at DESC"
		} else {
			orderStmt = "ORDER BY c.members_count ASC, c.created_at ASC"
		}
		if query.CursorSortValue != nil {
			cursorCount, err := strconv.Atoi(*query.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("failed to convert cursor sort value to int: %w", err)
			}
			args["cursor_members_count"] = cursorCount
			args["cursor_created_at"] = query.CursorCreatedAt
			if *query.Order == model.OrderDesc {
				conditions = append(conditions, "((c.members_count < @cursor_members_count) OR (c.members_count = @cursor_members_count AND c.created_at <= @cursor_created_at))")
			} else {
				conditions = append(conditions, "((c.members_count > @cursor_members_count) OR (c.members_count = @cursor_members_count AND c.created_at >= @cursor_created_at))")
			}
		}
	} else {
		if *query.Order == model.OrderDesc {
			orderStmt = "ORDER BY c.created_at DESC"
		} else {
			orderStmt = "ORDER BY c.created_at ASC"
		}
		if query.CursorCreatedAt != nil {
			args["cursor_created_at"] = *query.CursorCreatedAt
			if *query.Order == model.OrderDesc {
				conditions = append(conditions, "c.created_at <= @cursor_created_at")
			} else {
				conditions = append(conditions, "c.created_at >= @cursor_created_at")
			}
		}
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}
	stmt += " " + orderStmt

	limit := 20
	args["limit_plus_one"] = limit + 1
	stmt += " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to process get communities query: %w", err)
	}

	communities, err := pgx.CollectRows(rows, pgx.RowToStructByName[community.MiniCommunity])
	if err != nil {
		return nil, fmt.Errorf("Failed to collect rows: %w", err)
	}

	if len(communities) < limit+1 {
		var cursorSortValue string
		var cursorCreatedAt time.Time
		if query.CursorSortValue != nil {
			cursorSortValue = *query.CursorSortValue
		}
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[community.MiniCommunity]{
			Data:            communities,
			CursorSortValue: cursorSortValue,
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	if *query.Sort == model.SortByMembersCount {
		cursorSortValue := communities[limit].MembersCount
		cursorSortCreatedAt := communities[limit].CreatedAt

		return &model.CursorPaginatedResponse[community.MiniCommunity]{
			Data:            communities[:limit],
			CursorSortValue: strconv.Itoa(cursorSortValue),
			CursorCreatedAt: cursorSortCreatedAt,
			HasMore:         true,
		}, nil
	}

	cursorSortValue := communities[limit].CreatedAt
	cursorSortCreatedAt := communities[limit].CreatedAt

	return &model.CursorPaginatedResponse[community.MiniCommunity]{
		Data:            communities[:limit],
		CursorSortValue: cursorSortValue.Format(time.RFC3339Nano),
		CursorCreatedAt: cursorSortCreatedAt,
		HasMore:         true,
	}, nil

}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) JoinCommunity(ctx context.Context, userID uuid.UUID, communityID uuid.UUID) (*community.CommunityMember, error) {

	banned, err := r.IsBannedFromCommunity(ctx, userID, communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ban status for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if *banned {
		code := "USER_IS_BANNED"
		return nil, errs.NewBadRequestError("user is banned from this community", false, &code, nil, nil)
	}

	isMember, err := r.IsMember(ctx, communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check membership for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if *isMember {
		code := "ALREADY_MEMBER"
		return nil, errs.NewBadRequestError("already a member of this community", false, &code, nil, nil)
	}

	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		INSERT INTO community_members (user_id, community_id, role)
		VALUES (@user_id, @community_id, @role)
		RETURNING *
	`, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
		"role":         community.MemberRole,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to join community %s for user_id %s: %w", communityID, userID, err)
	}
	member, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[community.CommunityMember])
	if err != nil {
		return nil, fmt.Errorf("failed to parse new community member: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO community_follows (follower_id, community_id)
		VALUES (@follower_id, @community_id)
		ON CONFLICT (follower_id, community_id) DO NOTHING
	`, pgx.NamedArgs{
		"follower_id":  userID,
		"community_id": communityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create default follow on join: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &member, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) LeaveCommunity(ctx context.Context, userID uuid.UUID, communityID uuid.UUID) error {

	isAdmin, err := r.IsAdmin(ctx, communityID, userID)
	if err != nil {
		return fmt.Errorf("failed to check admin status for user_id %s and community_id %s: %w", userID, communityID, err)
	}
	if *isAdmin {
		code := "ADMIN_CANNOT_LEAVE"
		return errs.NewBadRequestError("admin must transfer ownership before leaving the community", false, &code, nil, nil)
	}

	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM community_members
		WHERE user_id = @user_id AND community_id = @community_id
	`, pgx.NamedArgs{
		"user_id":      userID,
		"community_id": communityID,
	})
	if err != nil {
		return fmt.Errorf("failed to leave community %s for user_id %s: %w", communityID, userID, err)
	}
	if result.RowsAffected() == 0 {
		code := "NOT_A_MEMBER"
		return errs.NewNotFoundError("not a member of this community", false, &code)
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
