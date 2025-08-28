package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/vrv501/simple-api/internal/constants"
)

func (m *mongoClient) performAdvisoryLockDBOperation(ctx context.Context, uniqueID bson.ObjectID,
	fn func(aCtx context.Context) (any, error)) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()
	lockExpiresAt := time.Now().Add(constants.DefaultTimeout).UTC()

	var (
		i   int
		err error
	)
	for i = range retriesForLease {
		_, err = m.mongoDbHandler.Collection(leasesCollection).UpdateOne(
			ctx,
			bson.M{
				iDField:          uniqueID,
				lockedUntilField: bson.Null{},
			},
			bson.M{setOperator: bson.M{lockedUntilField: lockExpiresAt}},
			options.UpdateOne().SetUpsert(true),
		)
		if err == nil {
			break
		}
		if !mongo.IsDuplicateKeyError(err) {
			return nil, err
		}

		time.Sleep(leaseWaitTime)
	}
	if i >= retriesForLease {
		return nil, errors.New("failed to acquire lock on uniqueID " + uniqueID.Hex())
	}
	defer func() {
		m.mongoDbHandler.Collection(leasesCollection).UpdateOne(
			ctx,
			bson.M{iDField: uniqueID},
			bson.M{setOperator: bson.M{lockedUntilField: bson.Null{}}})
	}()

	return fn(ctx)
}
