package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/muhammadnajie/kp/internal/resources"
	"strconv"

	"github.com/muhammadnajie/kp/graph/generated"
	"github.com/muhammadnajie/kp/graph/model"
	"github.com/muhammadnajie/kp/internal/auth"
	"github.com/muhammadnajie/kp/internal/links"
	"github.com/muhammadnajie/kp/internal/pkg/jwt"
	"github.com/muhammadnajie/kp/internal/users"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ExtractUserContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	var link links.Link
	link.User = user
	link.Title = input.Title
	link.Address = input.Address
	linkID := link.Save()

	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   link.Title,
		Address: link.Address,
		User:    graphqlUser,
	}, nil
}

func (r *mutationResolver) UpdateLink(ctx context.Context, input model.UpdateLink) (*model.Link, error) {
	user := auth.ExtractUserContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	ID, _ := strconv.Atoi(input.ID)
	userID, _ := strconv.Atoi(user.ID)
	link, err := links.GetByID(ID, userID)
	if err != nil {
		return nil, err
	}
	link.Title = input.Title
	link.Address = input.Address
	_, err = link.Update()
	if err != nil {
		return nil, err
	}
	return &model.Link{
		ID:      link.ID,
		Title:   link.Title,
		Address: link.Address,
	}, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, input model.DeleteLink) (string, error) {
	user := auth.ExtractUserContext(ctx)
	if user == nil {
		return "", fmt.Errorf("access denied")
	}
	var link links.Link
	link.ID = input.ID
	_, err := link.Delete()
	if err != nil {
		return "", nil
	}
	return "successfully delete", nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	valid, err := resources.ValidateUser(input.Username, input.Password)
	if !valid {
		return "", err
	}

	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	authenticated := user.Authenticate()
	if !authenticated {
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Links(ctx context.Context, title *string) ([]*model.Link, error) {
	user := auth.ExtractUserContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	var resultLinks []*model.Link
	var dbLinks []links.Link
	var err error
	if title != nil {
		dbLinks, err = links.GetByTitle(*title, user.ID)
	} else {
		dbLinks, err = links.GetAll(user.ID)
	}
	if err != nil {
		return nil, err
	}
	for _, link := range dbLinks {
		var u = &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{
			ID:      link.ID,
			Title:   link.Title,
			Address: link.Address,
			User:    u,
		})
	}
	return resultLinks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
