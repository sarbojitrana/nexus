package repository

import (
	"context"
	"errors"
	"fmt"
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
	stmt := `
		INSERT INTO communities(
			admin_id,
			name,
			slug,
			description,
			avatar_key,
			banner_key,
			can_post
		)
		VALUES(
			@admin_id,
			@name,
			@slug,
			@description,
			@avatar_key,
			@banner_key,
			@can_post
		)
		RETURNING
		*
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"admin_id":    payload.AdminID,
		"name":        payload.Name,
		"slug":        payload.Slug,
		"description": payload.Description,
		"avatar_key":  *payload.AvatarKey,
		"banner_key":  *payload.BannerKey,
		"can_post":    *payload.CanPost,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create community: %w", err)
	}

	community, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[community.Community])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to struct: %w", err)
	}

	return &community, nil
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

	if payload.CanPost != nil {
		args["can_post"] = *payload.CanPost
		setClauses = append(setClauses, "can_post = @can_post")
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("No fields provided to update for community_id %s: %w", communityID, err)
	}

	stmt += " " + strings.Join(setClauses, " AND ") + " WHERE community_id = @community_id"

	var community community.Community
	err = r.server.DB.Pool.QueryRow(ctx, stmt, args).Scan(&community)

	if err != nil {
		return nil, fmt.Errorf("Failed to update the community settings for community_id %s: %w", communityID, err)

	}
	return &community, nil
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
		return fmt.Errorf("community not found")
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) ChangeMemberRoleInCommunity(ctx context.Context, communityID uuid.UUID, userID uuid.UUID, payload *community.ChangeMemberRoleInCommunityPayload) (*community.CommunityMember, error) {

	check, err := r.IsAdmin(ctx, communityID, userID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was the admin for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if *check == false {
		return nil, fmt.Errorf("The user with user_id %s is not the admin of the community", userID)
	}

	stmt := `
		UPDATE community_members SET(
			role
		)
		VALUES(
			@role
		)
		WHERE community_id = @community_id AND AND user_id = @user_id 
	`
	var communityMember community.CommunityMember
	err = r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"role": payload.NewRole,
	}).Scan(&communityMember)

	if err != nil {
		return nil, fmt.Errorf("Failed to update the role of the member: %w", err)
	}

	return &communityMember, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) GetCommunityPostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, communityID uuid.UUID) (*post.PopulatedPost, error) {
	
	isUserBanned, err := r.IsBannedFromCommunity(ctx, userID, communityID)

	if err != nil{
		return nil, fmt.Errorf("Could Not check if the user was banned by the community: %w",err)
	}

	if *isUserBanned == true {
		code := "USER IS BANNED FROM COMMUNITY"
		return nil, errs.NewBadRequestError{
			"user is banned",
			false,
			&code,
		}
	}

	stmt := `
		SELECT p.* , 
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY
					pm.created_at DESC,
					pm.id DESC
			) FILTER(
				WHERE pm.id IS NOT NULL 
			),
			'[]' :: JSONB
		) AS post_media,
		json_build_object(
			'communityId', community.id
			'communityName', community.name,
			'communityAvatarKey', community.avatar_key
		) AS mini_community
		FROM posts p
		LEFT JOIN post_media pm ON pm.post_id = p.id
		LEFT JOIN communities community ON community.id = @community_id
		WHERE p.id = @post_id
		GROUP BY p.id
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"post_id":      postID,
		"community_id": communityID,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to process get post by id query for post_id %s: %w", postID, err)
	}

	post, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to the struct: %w", err)
	}

	return &post, nil

}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) GetCommunityMembers(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, query *community.GetCommunityMembersQuery) (*model.CursorPaginatedResponse[community.MiniCommunityUser], error) {
		
	isUserBanned, err := r.IsBannedFromCommunity(ctx, userID, communityID)

	if err != nil{
		return nil, fmt.Errorf("Could Not check if the user was banned by the community: %w",err)
	}

	if *isUserBanned == true {
		code := "USER IS BANNED FROM COMMUNITY"
		return nil, errs.NewBadRequestError{
			"user is banned",
			false,
			&code,
		}
	}	

	stmt := `
		SELECT cm.*,
		COALESCE(
			json_agg(
				to_jsonb(camel(u))
				ORDER BY
					created_at DESC,
					id DESC
			), FILTER(
				WHERE id IS NOT NULL
			),
			'[]' :: JSONB
		) AS users

		FROM community_members cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.community_id = @community_id

	`

	limit := 20

	args := pgx.NamedArgs{
		"community_id":   communityID,
		"limit_plus_one": limit + 1,
	}

	orderStmt := ""

	if *query.Order == model.OrderDesc {
		orderStmt = "ORDER BY joined_at DESC"
		if query.CursorSortValue != nil {
			stmt += " AND joined_at <= @joined_at "
			args["joined_at"] = query.CursorSortValue
		}

	} else {
		orderStmt += "ORDER BY joined_at ASC"
		if query.CursorSortValue != nil {
			stmt += " AND joined_at >= @joined_at "
			args["joined_at"] = query.CursorSortValue
		}
	}

	stmt += orderStmt + " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to process get community members of community with comunity_id %s: %w", communityID, err)
	}

	communityMembers, err := pgx.CollectRows(rows, pgx.RowToStructByName[community.MiniCommunityUser])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.CursorPaginatedResponse[community.MiniCommunityUser]{
				Data:            []community.MiniCommunityUser{},
				CursorSortValue: *query.CursorSortValue,
				CursorCreatedAt: *query.CursorCreatedAt,
				HasMore:         false,
			}, nil
		}
		return nil, fmt.Errorf("Failed to parse the rows to the struct: %w", err)
	}

	if len(communityMembers) < limit+1 {
		length := len(communityMembers)
		return &model.CursorPaginatedResponse[community.MiniCommunityUser]{
			Data:            communityMembers[:length],
			CursorSortValue: *query.CursorSortValue,
			CursorCreatedAt: *query.CursorCreatedAt,
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
		INSERT INTO
		community_reports(
			reporter_id,
			community_id,
			post_id,
			reason,
			status
		)
		VALUES(
			@reporter_id,
			@community_id,
			@post_id,
			@reason,
			@status
		)
		RETURNING
		*
	`
	var communityReport community.CommunityReport
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"reporter_id":  userID,
		"community_id": payload.CommunityID,
		"post_id":      payload.PostID,
		"reason":       payload.Reason,
		"status":       community.ReportPending,
	}).Scan(&communityReport)

	if err != nil {
		return nil, fmt.Errorf("Failed to create a report: %w", err)
	}
	return &communityReport, nil
}


//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) ResolveCommunityPostReport(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.ResolveCommunityPostReportPayload) (*community.CommunityReport, error){
	check, err := r.IsModerator(ctx, communityID, userID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was the moderator/admin for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if *check == false {
		return nil, fmt.Errorf("The user with user_id %s is not the moderator/admin of the community", userID)
	}

	stmt := `
		UPDATE community_reports SET
		status = @updated_status
		WHERE id = @report_id
		RETURNING *
	`
	var report community.CommunityReport
	err = r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"report_id" : payload.ReportID,
		"updated_status" : payload.UpdatedStatus,
	}).Scan(&report)

	if err != nil {
		return nil, fmt.Errorf("Failed to update the report with report_id %s: %w", payload.ReportID, err)
	}

	return &report, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) DeleteCommunityPost(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.DeleteCommunityPostPayload) error {
	check, err := r.IsModerator(ctx, communityID, userID)

	if err != nil {
		return fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if *check == false {
		return fmt.Errorf("The user with user_id %s is not a moderator/admin of the community", userID)
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
		return fmt.Errorf("Post with post_id %s not found or does not belong to community_id %s", payload.PostID, communityID)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) BanUserFromCommunity( ctx context.Context, userID uuid.UUID, communityID uuid.UUID, payload *community.BanCommunityMemberPayload ) (*community.BannedFromCommunityUser, error) {
	check, err := r.IsModerator(ctx, communityID, userID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if *check == false {
		return nil,fmt.Errorf("The user with user_id %s is not a moderator/admin of the community", userID)
	}

	stmt := `
		INSERT INTO banned_from_community_users (
			community_id,
			user_id,
			duration
		)
		VALUES (
			@community_id,
			@user_id_to_ban,
			@duration
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"community_id":   communityID,
		"user_id_to_ban": payload.UserIDToBan,
		"duration":       payload.Duration,
	})

	if err != nil {
		return nil, fmt.Errorf(
			"failed to ban user %s from community %s by user %s: %w",
			payload.UserIDToBan,
			communityID,
			userID,
			err,
		)
	}

	ban, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[community.BannedFromCommunityUser],
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse banned user for community %s: %w",
			communityID,
			err,
		)
	}

	return &ban, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetCommunityReports( ctx context.Context, userID uuid.UUID, communityID uuid.UUID, query *community.GetCommunityReportsQuery) (*model.CursorPaginatedResponse[community.CommunityReport], error) {

	check, err := r.IsModerator(ctx, communityID, userID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was a moderator for user_id %s and community_id %s: %w", userID, communityID, err)
	}

	if *check == false {
		return nil, fmt.Errorf("The user with user_id %s is not a moderator/admin of the community", userID)
	}

	stmt := `
		SELECT *
		FROM community_reports
		WHERE community_id = @community_id
	`
	limit := 20

	args := pgx.NamedArgs{
		"community_id": communityID,
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

	if len(conditions) > 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	stmt += `
		ORDER BY created_at DESC
		LIMIT @limit_plus_one
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to query community reports for community_id %s: %w",
			communityID,
			err,
		)
	}

	reports, err := pgx.CollectRows(rows, pgx.RowToStructByName[community.CommunityReport])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.CursorPaginatedResponse[community.CommunityReport]{
				Data:             []community.CommunityReport{},
				CursorSortValue:  "",
				CursorCreatedAt:  time.Time{},
				HasMore:          false,
			}, nil
		}

		return nil, fmt.Errorf(
			"failed to collect community reports for community_id %s: %w",
			communityID,
			err,
		)
	}

	if len(reports) < limit + 1 {
		return &model.CursorPaginatedResponse[community.CommunityReport]{
			Data:             reports,
			CursorSortValue:  "",
			CursorCreatedAt:  time.Time{},
			HasMore:          false,
		}, nil
	}

	return &model.CursorPaginatedResponse[community.CommunityReport]{
		Data: reports[:limit],
		CursorSortValue: reports[limit].CreatedAt.Format(time.RFC3339Nano),
		CursorCreatedAt: reports[limit].CreatedAt,
		HasMore: true,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) GetReportByID( ctx context.Context, userID uuid.UUID, communityID uuid.UUID, reportID uuid.UUID ) (*community.CommunityReport, error) {

	stmt := `
		SELECT *
		FROM community_reports
		WHERE id = @report_id
		AND community_id = @community_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"report_id":    reportID,
		"community_id": communityID,
	})

	if err != nil {
		return nil, fmt.Errorf(
			"failed to query report_id %s in community_id %s: %w",
			reportID,
			communityID,
			err,
		)
	}

	report, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[community.CommunityReport],
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse report_id %s: %w",
			reportID,
			err,
		)
	}

	return &report, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *CommunityRepository) IsBannedFromCommunity( ctx context.Context, userID uuid.UUID, communityID uuid.UUID,) (*bool, error) {

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
