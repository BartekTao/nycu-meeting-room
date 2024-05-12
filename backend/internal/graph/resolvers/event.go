package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/internal/common"
	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/middleware"
)

// Room is the resolver for the room field.
func (r *eventResolver) Room(ctx context.Context, obj *domain.Event) (*domain.Room, error) {
	if obj.RoomID == nil {
		return nil, nil
	}

	room, err := r.roomService.GetByID(ctx, *obj.RoomID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// Participants is the resolver for the participants field.
func (r *eventResolver) Participants(ctx context.Context, obj *domain.Event) ([]domain.User, error) {
	return r.userService.GetByIDs(ctx, obj.ParticipantsIDs)
}

// Creator is the resolver for the creator field.
func (r *eventResolver) Creator(ctx context.Context, obj *domain.Event) (*domain.User, error) {
	return r.userService.GetByID(ctx, obj.CreatorID)
}

// UpsertEvent is the resolver for the upsertEvent field.
func (r *mutationResolver) UpsertEvent(ctx context.Context, input model.UpsertEventInput) (*domain.Event, error) {
	claims, _ := ctx.Value(middleware.UserCtxKey).(middleware.MeetingCenterClaims)

	upsertEvent := app.UpsertEventRequest{
		ID:              input.ID,
		Title:           input.Title,
		Description:     input.Description,
		StartAt:         input.StartAt,
		EndAt:           input.EndAt,
		RoomID:          input.RoomID,
		ParticipantsIDs: input.ParticipantsIDs,
		Notes:           input.Notes,
		RemindAt:        input.RemindAt,
		UpdaterID:       claims.Sub,
	}

	event, upsert_err := r.eventService.Upsert(ctx, upsertEvent)
	if upsert_err != nil {
		return nil, upsert_err
	}

	return event, nil
}

// DeleteEvent is the resolver for the deleteEvent field.
func (r *mutationResolver) DeleteEvent(ctx context.Context, id string) (*domain.Event, error) {
	event, err := r.eventService.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// UserEvent is the resolver for the userEvent field.
func (r *queryResolver) UserEvent(ctx context.Context, userID string) ([]domain.Event, error) {
	panic(fmt.Errorf("not implemented: UserEvent - userEvent"))
}

// Event is the resolver for the event field.
func (r *queryResolver) Event(ctx context.Context, id string) (*domain.Event, error) {
	event, err := r.eventService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// PaginatedAvailableRooms is the resolver for the paginatedAvailableRooms field.
func (r *queryResolver) PaginatedAvailableRooms(ctx context.Context, startAt int64, endAt int64, first *int, after *string) (*model.RoomConnection, error) {
	panic(fmt.Errorf("not implemented: PaginatedAvailableRooms - paginatedAvailableRooms"))
}

// UserEvents is the resolver for the userEvents field.
func (r *queryResolver) UserEvents(ctx context.Context, userIDs []string, startAt int64, endAt int64) ([]*model.UserEvent, error) {
	userEvents, err := r.eventService.GetUserEvents(ctx, userIDs, startAt, endAt)
	if err != nil {
		return nil, err
	}

	res := make([]*model.UserEvent, len(userEvents))
	for userID, userEvent := range userEvents {
		res = append(res, &model.UserEvent{
			User: &domain.User{
				ID: common.ToPtr(userID),
			},
			Events: userEvent,
		})
	}
	return res, nil
}

// User is the resolver for the user field.
func (r *userEventResolver) User(ctx context.Context, obj *model.UserEvent) (*domain.User, error) {
	return r.userService.GetByID(ctx, *obj.User.ID)
}

// Event returns graph.EventResolver implementation.
func (r *Resolver) Event() graph.EventResolver { return &eventResolver{r} }

// UserEvent returns graph.UserEventResolver implementation.
func (r *Resolver) UserEvent() graph.UserEventResolver { return &userEventResolver{r} }

type eventResolver struct{ *Resolver }
type userEventResolver struct{ *Resolver }