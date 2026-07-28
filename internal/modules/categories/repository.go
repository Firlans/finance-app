package categories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Save(ctx context.Context, category *Category) error
	GetCategories(ctx context.Context, userID *string, typeCategory string) ([]Category, error)
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) Save(ctx context.Context, category *Category) error {
	typeCategory := category.TypeCategory
	if typeCategory == "" {
		typeCategory = "expense"
	}
	query := `INSERT INTO categories (name, description, type_category, user_id, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id int
	err := r.QueryRow(ctx, query, category.Name, category.Description, typeCategory, category.UserID).Scan(&id)
	if err != nil {
		return err
	}

	category.ID = id
	category.TypeCategory = typeCategory
	return nil
}

func (r *repository) GetCategories(ctx context.Context, userID *string, typeCategory string) ([]Category, error) {
	query := `SELECT id, name, description, COALESCE(type_category, 'expense'), created_at FROM categories WHERE deleted_at IS NULL AND (user_id = $1 OR user_id IS NULL)`
	var args []interface{}
	args = append(args, *userID)

	if typeCategory != "" {
		query += ` AND (type_category = $2 OR type_category = 'both')`
		args = append(args, typeCategory)
	}

	rows, err := r.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories = make([]Category, 0)
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.TypeCategory, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *repository) GetCategoryByID(ctx context.Context, id int) (*Category, error) {
	query := `SELECT id, name, description, COALESCE(type_category, 'expense'), created_at FROM categories WHERE id = $1 AND deleted_at IS NULL`
	row := r.QueryRow(ctx, query, id)

	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.Description, &category.TypeCategory, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *repository) UpdateCategory(ctx context.Context, category *Category) error {
	query := `UPDATE categories SET name = $1, description = $2, type_category = COALESCE(NULLIF($3, ''), type_category) WHERE id = $4 AND deleted_at IS NULL`
	_, err := r.Exec(ctx, query, category.Name, category.Description, category.TypeCategory, category.ID)
	return err
}

func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	query := `UPDATE categories SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.Exec(ctx, query, id)
	return err
}
