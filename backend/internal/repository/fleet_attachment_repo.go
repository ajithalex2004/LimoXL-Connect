package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)


type PostgresFleetAttachmentRepo struct {
	DB *sql.DB
}

func NewPostgresFleetAttachmentRepo(db *sql.DB) *PostgresFleetAttachmentRepo {
	return &PostgresFleetAttachmentRepo{DB: db}
}

func (r *PostgresFleetAttachmentRepo) Create(ctx context.Context, att *models.FleetAttachment) error {
	query := `
		INSERT INTO fleet_attachments (
			id, entity_id, entity_type, file_name, file_url, file_type, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	att.ID = uuid.New()
	att.CreatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		att.ID, att.EntityID, att.EntityType, att.FileName, att.FileURL, att.FileType, att.CreatedAt,
	)
	return err
}

func (r *PostgresFleetAttachmentRepo) ListByEntity(ctx context.Context, entityID uuid.UUID, entityType string) ([]models.FleetAttachment, error) {
	query := `
		SELECT id, entity_id, entity_type, file_name, file_url, file_type, created_at
		FROM fleet_attachments
		WHERE entity_id = $1 AND entity_type = $2
		ORDER BY created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, entityID, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.FleetAttachment
	for rows.Next() {
		var att models.FleetAttachment
		if err := rows.Scan(
			&att.ID, &att.EntityID, &att.EntityType, &att.FileName, &att.FileURL, &att.FileType, &att.CreatedAt,
		); err != nil {
			return nil, err
		}
		attachments = append(attachments, att)
	}
	return attachments, nil
}

func (r *PostgresFleetAttachmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM fleet_attachments WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
