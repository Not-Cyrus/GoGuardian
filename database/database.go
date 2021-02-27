package database

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *MongoDB) CreateGuild(botUser *discordgo.User, guild *discordgo.Guild) {
	_, err = db.FindData(guild.ID)
	if err == nil {
		return
	}

	if _, err = db.Collection.InsertOne(context.Background(), bson.M{"users": bson.A{guild.OwnerID}, "log-channel": "nil", "prefix": ">", "guild_id": guild.ID, "anti-invite": "off"}); err != nil {
		return
	}

	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"guild_id": guild.ID}, bson.M{"$push": bson.M{"users": botUser.ID}}); err != nil {
		return
	}
	//create fucking stuff
}

func (db *MongoDB) DeleteGuild(guildID string) bool {
	if _, err = db.Collection.DeleteOne(context.Background(), bson.M{"guild_id": guildID}); err != nil {
		return false
	}
	return true
}

func (db *MongoDB) FindData(guildID string) (bson.M, error) {
	var data bson.M
	if err = db.Collection.FindOne(context.Background(), bson.M{"guild_id": guildID}).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (db *MongoDB) IsWhitelisted(guildID string, userID string) bool {
	var data bson.M

	data, err = db.FindData(guildID)
	if err != nil {
		return false
	}

	for _, whitelistedID := range data["users"].(bson.A) {
		if whitelistedID == userID {
			return true
		}
	}
	return false
}

func (db *MongoDB) SetData(guildID string, index string, value string) (bool, error) {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"guild_id": guildID}, bson.M{"$set": bson.M{index: value}}, &options.UpdateOptions{}); err != nil {
		return false, err
	}
	return true, nil
}

func (db *MongoDB) SetWhitelist(guildID string, member *discordgo.User, whitelist bool) (bool, error) {
	whitelisted := db.IsWhitelisted(guildID, member.ID)

	if whitelist && whitelisted {
		return false, fmt.Errorf("I couldn't seem to change their whitelist status, Maybe check whitelists?")
	}

	switch whitelist {
	case true:
		if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"guild_id": guildID}, bson.M{"$push": bson.M{"users": member.ID}}, &options.UpdateOptions{}); err != nil {
			return false, err
		}
	case false:
		if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"guild_id": guildID}, bson.M{"$pull": bson.M{"users": member.ID}}, &options.UpdateOptions{}); err != nil {
			return false, err
		}
	}

	return true, nil
}

func SetupDB() MongoDB {
	var db = MongoDB{}

	db.Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/Bot?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}

	err = db.Client.Connect(db.Ctx)
	if err != nil {
		panic(err)
	}

	db.Database = db.Client.Database("Bot")
	db.Collection = db.Database.Collection("whitelist")

	return db
}

var (
	cancel   func()
	Database = SetupDB()
	err      error
)

type (
	MongoDB struct {
		Collection *mongo.Collection
		Client     *mongo.Client
		Ctx        context.Context
		Database   *mongo.Database
	}
)
