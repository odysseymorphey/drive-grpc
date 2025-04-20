package postgres

import (
	"context"
	"database/sql"
	"drive-grpc/internal/server/models"
	"fmt"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) GetFileInfoByName(ctx context.Context, filename string, fileId string) (*models.FileInfo, error) {
	const op = "storage.postgres.GetFileInfoByName"

	const query = `SELECT id, file_name, file_path, creation_date, modification_date
				   FROM files
			  	   WHERE file_name=$1 AND id=$2`

	row := d.db.QueryRowContext(ctx, query, filename, fileId)

	var fInfo models.FileInfo
	err := row.Scan(
		&fInfo.ID,
		&fInfo.Filename,
		&fInfo.FilePath,
		&fInfo.CreationDate,
		&fInfo.ModificationDate)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %v", op, err)
	}

	return &fInfo, nil
}

func (d *Database) GetFilesInfoList(ctx context.Context) ([]models.FileInfo, error) {
	const op = "storage.postgres.GetFilesInfoList"

	const query = `SELECT id, file_name, file_path, creation_date, modification_date
				   FROM files`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Printf("%s: failed to execute query: %v", op, err)
		return nil, err
	}

	var fInfoList []models.FileInfo
	for rows.Next() {
		var fInfo models.FileInfo
		err := rows.Scan(
			&fInfo.ID,
			&fInfo.Filename,
			&fInfo.FilePath,
			&fInfo.CreationDate,
			&fInfo.ModificationDate)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %v", op, err)
		}

		fInfoList = append(fInfoList, fInfo)
	}

	return fInfoList, nil
}

func (d *Database) WriteFileInfo(ctx context.Context, info *models.FileInfo) error {
	const op = "storage.postgres.WriteFileInfo"

	const query = `INSERT INTO files
    			   (id, file_name, file_path, creation_date, modification_date)
				   VALUES ($1, $2, $3, $4, $5)`

	_, err := d.db.ExecContext(ctx, query,
		info.ID,
		info.Filename,
		info.FilePath,
		info.CreationDate,
		info.ModificationDate)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %v", op, err)
	}

	return nil
}
