package managers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/bahrom656/crud/cmd/app/middleware"
	"github.com/bahrom656/crud/pkg/customers"
	"github.com/bahrom656/crud/pkg/security"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//Service
type Service struct {
	pool *pgxpool.Pool
}

//NewService
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Manager struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Roles    []string  `json:"roles"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Qty   int64  `json:"qty"`
}

type Registration struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Service) Register(ctx context.Context, registration *Registration) (*Manager, error) {
	item := &Manager{}
	var id int64

	idManager, err := middleware.Authentication(ctx)
	if err != nil {
	}

	err = s.pool.QueryRow(ctx, `SELECT id FROM managers WHERE roles[0-9]='ADMIN' AND id=$1`, idManager).Scan(&id)
	if err != nil {
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
	}

	err = s.pool.QueryRow(ctx,
		`INSERT INTO managers(name, phone, password, roles) 
			VALUES ($1, $2, $3, $4) ON CONFLICT (phone) 
			DO NOTHING RETURNING id, name, phone, roles, active, created`,
		registration.Name, registration.Phone, hash, registration.Roles,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Roles,
		&item.Active,
		&item.Created,
	)
	if err == pgx.ErrNoRows {
	}
	if err != nil {
	}

	return item, nil
}

func (s *Service) AuthenticateManager(ctx context.Context, token string) (id int64, err error) {
	expireTime := time.Now()

	err = s.pool.QueryRow(
		ctx,
		`SELECT manager_id, expire FROM managers_tokens WHERE token = $1`,
		token,
	).Scan(&id, &expireTime)

	if err == pgx.ErrNoRows {
		return 0, security.ErrNoSuchUser
	}

	if err != nil {
		return 0, security.ErrInternal
	}

	if time.Now().After(expireTime) {
		return 0, security.ErrTokenExpired
	}
	return id, nil
}

func (s *Service) Token(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64

	err = s.pool.QueryRow(ctx,
		`SELECT id, password FROM managers WHERE phone=$1`, phone).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
	}
	if err != nil {
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx,
		`INSERT INTO managers_tokens(token, manager_id) VALUES ($1, $2)`, token, id)
	if err != nil {
	}

	return token, nil

}

func (s *Service) IDByToken(ctx context.Context, token string) (id int64, err error) {

	err = s.pool.QueryRow(ctx,
		`SELECT manager_id FROM managerS_tokens WHERE token=$1`, token).Scan(id)
	if err != nil {

	}

	return id, nil
}

func (s *Service) GetProducts(ctx context.Context) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, price, qty FROM products WHERE active ORDER BY id LIMIT 500`)
	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}
	if err != nil {
		return nil, security.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &Product{}
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Qty)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) RemoveProductsByID(ctx context.Context, id int64) (err error) {
	_, err = s.pool.Query(ctx,
		`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
	}

	return nil
}

func (s *Service) RemoveCustomerByID(ctx context.Context, id int64) (err error) {
	_, err = s.pool.Query(ctx,
		`DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
	}

	return nil
}

func (s *Service) GetCustomer(ctx context.Context) ([]*customers.Customer, error) {
	items := make([]*customers.Customer, 0)
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, phone, active, created FROM customers WHERE active = true ORDER BY id LIMIT 500`)
	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}
	if err != nil {
		return nil, security.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &customers.Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) ManagerChangeCustomer(ctx context.Context, customer *customers.Customer) (*customers.Customer, error) {
	item := &customers.Customer{}

	err := s.pool.QueryRow(ctx,
		`UPDATE customers 
			SET name = $1, phone = $2, active = $3 
			WHERE id=$4 
			RETURNING id, name, phone, active`,
		customer.Active, customer.Name, customer.Phone, customer.ID).Scan(&item.ID, &item.Name, &item.Phone, &item.Active)
	if err != nil {
	}
	return item, nil
}

func (s *Service) ManagerChangeProduct(ctx context.Context, product *Product) (*Product, error) {
	item := &Product{}

	err := s.pool.QueryRow(ctx,
		`UPDATE products 
			SET name = $1, price = $2, qty = $3 
			WHERE id=$4 
			RETURNING id, name, price, qty`,
		product.Name, product.Price, product.Qty, product.ID).Scan(&item.ID, &item.Name, &item.Price, &item.Qty)
	if err != nil {
	}
	return item, nil

}
