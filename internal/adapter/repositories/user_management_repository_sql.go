package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	sqlc "music-streaming/internal/adapter/sql/sqlc"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

type SQLUserManagementRepository struct {
	queries *sqlc.Queries
	db      *pgx.Conn
	redis   *redis.Client
}

func NewSQLUserManagementRepository(db *pgx.Conn, redisClient *redis.Client) *SQLUserManagementRepository {
	return &SQLUserManagementRepository{
		queries: sqlc.New(db),
		db:      db,
		redis:   redisClient,
	}
}

const (
	userCacheKeyPrefix = "user:"
	userCacheTTL       = 5 * time.Minute
)

func (r *SQLUserManagementRepository) getUserCacheKey(username string) string {
	return userCacheKeyPrefix + username
}

func (r *SQLUserManagementRepository) invalidateUserCache(ctx context.Context, username string) error {
	return r.redis.Del(ctx, r.getUserCacheKey(username)).Err()
}

func (r *SQLUserManagementRepository) getUserFromCache(ctx context.Context, username string) (domain.User, bool) {
	key := r.getUserCacheKey(username)
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, false
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return domain.User{}, false
	}

	return user, true
}

func (r *SQLUserManagementRepository) setUserInCache(ctx context.Context, user domain.User) error {
	key := r.getUserCacheKey(user.Username)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, val, userCacheTTL).Err()
}

func (r *SQLUserManagementRepository) CreateUser(ctx context.Context, user domain.User) error {
	// Check if user already exists
	_, err := r.queries.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return &ports.FailedOperationError{Description: "User already exists"}
	}
	if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Determine if admin user
	var sqlUser sqlc.User
	if user.AdminRole {
		sqlUser, err = r.queries.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
	} else {
		sqlUser, err = r.queries.CreateDefaultUser(ctx, sqlc.CreateDefaultUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update all other fields if they were provided
	domainUser := toDomainUser(sqlUser)
	// Update fields from the provided user
	domainUser.LdapAuthenticated = user.LdapAuthenticated
	domainUser.AdminRole = user.AdminRole
	domainUser.SettingsRole = user.SettingsRole
	domainUser.StreamRole = user.StreamRole
	domainUser.JukeboxRole = user.JukeboxRole
	domainUser.DownloadRole = user.DownloadRole
	domainUser.UploadRole = user.UploadRole
	domainUser.PlaylistRole = user.PlaylistRole
	domainUser.CoverArtRole = user.CoverArtRole
	domainUser.CommentRole = user.CommentRole
	domainUser.PodcastRole = user.PodcastRole
	domainUser.ShareRole = user.ShareRole
	domainUser.VideoConversionRole = user.VideoConversionRole
	domainUser.MusicfolderId = user.MusicfolderId
	domainUser.MaxBitRate = user.MaxBitRate

	if err := r.UpdateUser(ctx, user.Username, domainUser); err != nil {
		return fmt.Errorf("failed to update user fields: %w", err)
	}

	// Cache the created user
	_ = r.setUserInCache(ctx, domainUser)

	return nil
}

func (r *SQLUserManagementRepository) UpdateUser(ctx context.Context, username string, user domain.User) error {
	// First, check if user exists
	_, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Convert musicFolderId array to comma-separated string
	musicFolderIdStr := ""
	if len(user.MusicfolderId) > 0 {
		musicFolderIdStr = strings.Join(user.MusicfolderId, ",")
	}

	var musicFolderId pgtype.Text
	if musicFolderIdStr != "" {
		musicFolderId = pgtype.Text{String: musicFolderIdStr, Valid: true}
	}

	// Note: We need to check if UpdateUser query exists, if not we'll need to add it
	// For now, let's assume it will be generated after we add it to the SQL file
	sqlUser, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		Username:            username,
		Password:            user.Password,
		Email:               user.Email,
		Scrobblingenabled:   user.ScrobblingEnabled,
		Ldapauthenticated:   user.LdapAuthenticated,
		Adminrole:           user.AdminRole,
		Settingsrole:        user.SettingsRole,
		Streamrole:          user.StreamRole,
		Jukeboxrole:         user.JukeboxRole,
		Downloadrole:        user.DownloadRole,
		Uploadrole:          user.UploadRole,
		Playlistrole:        user.PlaylistRole,
		Coverartrole:        user.CoverArtRole,
		Commentrole:         user.CommentRole,
		Podcastrole:         user.PodcastRole,
		Sharerole:           user.ShareRole,
		Videoconversionrole: user.VideoConversionRole,
		Musicfolderid:       musicFolderId,
		Maxbitrate:          user.MaxBitRate,
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	_ = r.invalidateUserCache(ctx, username)

	// Update cache with new user data
	domainUser := toDomainUser(sqlUser)
	_ = r.setUserInCache(ctx, domainUser)

	return nil
}

func (r *SQLUserManagementRepository) DeleteUser(ctx context.Context, username string) error {
	_, err := r.queries.DeleteUser(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate cache
	_ = r.invalidateUserCache(ctx, username)

	return nil
}

func (r *SQLUserManagementRepository) GetUser(ctx context.Context, username string) (domain.User, error) {
	// Try to get from cache first
	if user, found := r.getUserFromCache(ctx, username); found {
		return user, nil
	}

	// If not in cache, get from database
	sqlUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	domainUser := toDomainUser(sqlUser)

	// Cache the user
	_ = r.setUserInCache(ctx, domainUser)

	return domainUser, nil
}

func (r *SQLUserManagementRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	sqlUsers, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	users := make([]domain.User, 0, len(sqlUsers))
	for _, sqlUser := range sqlUsers {
		domainUser := toDomainUser(sqlUser)
		users = append(users, domainUser)
	}

	return users, nil
}

// Helper function to convert SQL user to domain user
func toDomainUser(sqlUser sqlc.User) domain.User {
	user := domain.User{
		Username:            sqlUser.Username,
		Password:            sqlUser.Password,
		Email:               sqlUser.Email,
		ScrobblingEnabled:   sqlUser.Scrobblingenabled,
		LdapAuthenticated:   sqlUser.Ldapauthenticated,
		AdminRole:           sqlUser.Adminrole,
		SettingsRole:        sqlUser.Settingsrole,
		StreamRole:          sqlUser.Streamrole,
		JukeboxRole:         sqlUser.Jukeboxrole,
		DownloadRole:        sqlUser.Downloadrole,
		UploadRole:          sqlUser.Uploadrole,
		PlaylistRole:        sqlUser.Playlistrole,
		CoverArtRole:        sqlUser.Coverartrole,
		CommentRole:         sqlUser.Commentrole,
		PodcastRole:         sqlUser.Podcastrole,
		ShareRole:           sqlUser.Sharerole,
		VideoConversionRole: sqlUser.Videoconversionrole,
		MaxBitRate:          sqlUser.Maxbitrate,
	}

	// Convert musicFolderId from comma-separated string to array
	if sqlUser.Musicfolderid.Valid && sqlUser.Musicfolderid.String != "" {
		user.MusicfolderId = strings.Split(sqlUser.Musicfolderid.String, ",")
	} else {
		user.MusicfolderId = []string{}
	}

	return user
}
