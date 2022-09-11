package model

import "time"

type Flip struct {
	SellUUID interface{} `json:"SellUuid"`
	Flipper  string      `json:"Flipper"`
	Profit   int         `json:"Profit"`
	Sell     struct {
		UUID             string      `json:"uuid"`
		Count            int         `json:"count"`
		StartingBid      int         `json:"startingBid"`
		Tag              string      `json:"tag"`
		ItemName         interface{} `json:"itemName"`
		Start            time.Time   `json:"start"`
		End              time.Time   `json:"end"`
		AuctioneerID     string      `json:"auctioneerId"`
		ProfileID        interface{} `json:"profileId"`
		CoopMembers      interface{} `json:"CoopMembers"`
		HighestBidAmount int         `json:"highestBidAmount"`
		Bids             []struct {
			Bidder    string    `json:"bidder"`
			ProfileID string    `json:"profileId"`
			Amount    int       `json:"amount"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"bids"`
		AnvilUses    int           `json:"anvilUses"`
		Enchantments []interface{} `json:"enchantments"`
		NbtData      struct {
			Data struct {
				PetInfo string `json:"petInfo"`
				UID     string `json:"uid"`
			} `json:"Data"`
		} `json:"nbtData"`
		ItemCreatedAt time.Time `json:"itemCreatedAt"`
		Reforge       string    `json:"reforge"`
		Category      string    `json:"category"`
		Tier          string    `json:"tier"`
		Bin           bool      `json:"bin"`
		FlatNbt       struct {
			Type           string `json:"type"`
			Active         string `json:"active"`
			Exp            string `json:"exp"`
			Tier           string `json:"tier"`
			HideInfo       string `json:"hideInfo"`
			HeldItem       string `json:"heldItem"`
			CandyUsed      string `json:"candyUsed"`
			UUID           string `json:"uuid"`
			HideRightClick string `json:"hideRightClick"`
			UID            string `json:"uid"`
		} `json:"flatNbt"`
	} `json:"Sell"`
	Buy struct {
		Reforge          string        `json:"reforge"`
		Category         string        `json:"category"`
		Tier             string        `json:"tier"`
		Enchantments     []interface{} `json:"enchantments"`
		UUID             string        `json:"uuid"`
		Count            int           `json:"count"`
		StartingBid      int           `json:"startingBid"`
		Tag              string        `json:"tag"`
		ItemName         string        `json:"itemName"`
		Start            string        `json:"start"`
		End              string        `json:"end"`
		AuctioneerID     string        `json:"auctioneerId"`
		ProfileID        string        `json:"profileId"`
		Coop             interface{}   `json:"coop"`
		CoopMembers      interface{}   `json:"coopMembers"`
		HighestBidAmount int           `json:"highestBidAmount"`
		Bids             []struct {
			Bidder    string `json:"bidder"`
			ProfileID string `json:"profileId"`
			Amount    int    `json:"amount"`
			Timestamp string `json:"timestamp"`
		} `json:"bids"`
		NbtData struct {
			Data struct {
				PetInfo string `json:"petInfo"`
				UID     string `json:"uid"`
			} `json:"data"`
		} `json:"nbtData"`
		ItemCreatedAt string `json:"itemCreatedAt"`
		Bin           bool   `json:"bin"`
		FlatNbt       struct {
			Type      string `json:"type"`
			Active    string `json:"active"`
			Exp       string `json:"exp"`
			Tier      string `json:"tier"`
			HideInfo  string `json:"hideInfo"`
			HeldItem  string `json:"heldItem"`
			CandyUsed string `json:"candyUsed"`
			UUID      string `json:"uuid"`
			UID       string `json:"uid"`
		} `json:"flatNbt"`
	} `json:"Buy"`
	Finder                         interface{} `json:"Finder"`
	AmountOfFlipsFromBuyerOfTheDay int         `json:"amountOfFlipsFromBuyerOfTheDay"`
}
