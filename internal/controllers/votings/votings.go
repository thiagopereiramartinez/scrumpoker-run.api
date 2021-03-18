package votings

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/rooms"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/votings"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/utils"
	"time"
)

var ctx context.Context

// @Summary Vote in a room
// @Tags Votings
// @Param body body votings.RegisterVoteRequest true "Vote fields"
// @Param pincode path string true "Pin Code"
// @Accept json
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /votings/{pincode} [post]
func registerVote(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	body := new(votings.RegisterVoteRequest)
	if err := c.BodyParser(body); err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}

	if err := body.Validate(); err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	db := new(firestore.Client)
	container.Make(&db)

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	var room rooms.Room
	_ = roomSnap.DataTo(&room)
	room.Id = roomSnap.Ref.ID

	userSnap, err := db.Collection("rooms").Doc(room.Id).Collection("players").Doc(body.PlayerId).Get(ctx)
	if err != nil || !userSnap.Exists() {
		_ = utils.SendError(c, 404, errors.New("player not found"))
		return nil
	}

	if _, err := db.Collection("rooms").Doc(room.Id).Collection("votings").Doc(body.PlayerId).Set(ctx, map[string]interface{}{
		"vote":      body.Value,
		"timestamp": firestore.ServerTimestamp,
	}, firestore.MergeAll); err != nil {
		_ = utils.SendError(c, 500, errors.New("an error occurred while registering the vote"))
		return nil
	}

	return c.SendStatus(200)
}

// @Summary Reset votes in a room
// @Tags Votings
// @Param pincode path string true "Pin Code"
// @Accept json
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /votings/{pincode}/reset [delete]
func resetVotes(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	db := new(firestore.Client)
	container.Make(&db)

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	batch := db.Batch()
	docs, _ := db.Collection("rooms").Doc(roomSnap.Ref.ID).Collection("votings").Documents(ctx).GetAll()
	for _, snap := range docs {
		batch.Delete(snap.Ref)
	}
	_, _ = batch.Commit(ctx)

	return c.SendStatus(200)
}

// @Summary Get votes in a room
// @Tags Votings
// @Param pincode path string true "Pin Code"
// @Accept json
// @Produce json
// @Success 200 {array} votings.Vote
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /votings/{pincode} [get]
func getVotes(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	db := new(firestore.Client)
	container.Make(&db)

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	var votes = make([]votings.Vote, 0)
	vts, err := db.Collection("rooms").Doc(roomSnap.Ref.ID).Collection("votings").Documents(ctx).GetAll()
	for _, vote := range vts {

		var v = votings.Vote{
			PlayerId: vote.Ref.ID,
			Value: vote.Data()["vote"].(float64),
			VotedAt: vote.Data()["timestamp"].(time.Time),
		}
		pls, _ := db.Collection("rooms").Doc(roomSnap.Ref.ID).Collection("players").Doc(vote.Ref.ID).Get(ctx)
		v.PlayerName = pls.Data()["name"].(string)

		votes = append(votes, v)
	}

	return c.JSON(votes)
}

// Register endpoints
func Register(router fiber.Router) {

	ctx = context.Background()

	v := router.Group("/votings")
	v.Get(":pincode", getVotes)
	v.Post(":pincode", registerVote)
	v.Delete(":pincode/reset", resetVotes)

}
