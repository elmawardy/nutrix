package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Seeder struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

func (s *Seeder) SeedProducts() error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Check if the product with name "ProductSeeded" exists in the db
	var product models.Product
	err = client.Database("waha").Collection("recipes").FindOne(ctx, bson.M{"name": "ProductSeeded"}).Decode(&product)
	if err == nil {
		confirmed, err := s.Prompter.Confirmation("products already exists, do you want to insert new documents beside the current ones ?")
		if err != nil {
			return err
		}

		if !confirmed {
			return nil
		}
	}

	// Connected successfully

	// Get the material with name Motzarilla from the DB
	var material models.Material
	err = client.Database("waha").Collection("components").FindOne(ctx, bson.M{"name": "MotzarillaSeeded"}).Decode(&material)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			confirmation, err := s.Prompter.Confirmation("no seeded materials found, would you like to create them?")
			if err != nil {
				return err
			}

			if confirmation {
				err = s.SeedMaterials(true)
				if err != nil {
					return err
				}

				err = client.Database("waha").Collection("components").FindOne(ctx, bson.M{"name": "MotzarillaSeeded"}).Decode(&material)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}

	}

	material.Quantity = 15

	products := []models.Product{
		{
			Name:  "ProductSeeded",
			Price: 100.0,
			Materials: []models.Material{
				material,
			},
			SubProducts: []models.Product{
				{
					Name:  "Product 2",
					Price: 100.0,
					Materials: []models.Material{
						material,
					},
					SubProducts: []models.Product{},
				},
			},
		},
	}

	newValue := make([]interface{}, len(products))

	for i := range products {
		newValue[i] = products[i]
	}

	// Insert the products
	_, err = client.Database("waha").Collection("recipes").InsertMany(ctx, newValue)
	if err != nil {
		return err
	}

	s.Logger.Info("products seeded successfully !")

	return nil
}

func (s *Seeder) SeedCategories() error {

	categories := []models.Category{
		{
			Name:    "CategorySeeded",
			Recipes: []string{},
		},
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}
	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Set the database and collection
	db := client.Database("waha")
	collection := db.Collection("categories")

	// Count the number of documents in the categories collection
	count, err := collection.CountDocuments(ctx, bson.D{})

	if count > 0 {
		confirmed, err := s.Prompter.Confirmation("categories already exists, do you want to insert new documents beside the current ones ?")
		if err != nil {
			return err
		}

		if !confirmed {
			return nil
		}

		var product models.Product
		err = client.Database("waha").Collection("recipes").FindOne(ctx, bson.M{"name": "ProductSeeded"}).Decode(&product)

		if err == mongo.ErrNoDocuments {
			confirm, err := s.Prompter.Confirmation("seeded products not found, would you like to create one?")

			if err != nil {
				return err
			}

			if confirm {
				err = s.SeedProducts()
				if err != nil {
					return err
				}

				err = client.Database("waha").Collection("recipes").FindOne(ctx, bson.M{"name": "ProductSeeded"}).Decode(&product)
				if err != nil {
					return err
				}

				for index := range categories {
					categories[index].Recipes = append(categories[index].Recipes, product.Id)
				}
			}
		} else if err != nil {
			return err
		} else {
			for index := range categories {
				categories[index].Recipes = append(categories[index].Recipes, product.Id)
			}
		}

		newValue := make([]interface{}, len(categories))

		for i := range categories {
			newValue[i] = categories[i]
		}

		_, err = collection.InsertMany(ctx, newValue)
		if err != nil {
			return err
		}

		s.Logger.Info("categories seeded successfully!")

		return nil
	} else if err == mongo.ErrNoDocuments || err == nil {

		var product models.Product
		err = client.Database("waha").Collection("recipes").FindOne(ctx, bson.M{"name": "ProductSeeded"}).Decode(&product)

		if err == mongo.ErrNoDocuments {
			confirm, err := s.Prompter.Confirmation("No seeded products found, would you like to create one?")

			if err != nil {
				return err
			}

			if confirm {
				err = s.SeedProducts()
				if err != nil {
					return err
				}

				err = client.Database("waha").Collection("recipes").FindOne(ctx, bson.M{"name": "ProductSeeded"}).Decode(&product)
				if err != nil {
					return err
				}

				for index := range categories {
					categories[index].Recipes = append(categories[index].Recipes, product.Id)
				}
			}

		} else if err != nil {
			return err
		} else {
			for index := range categories {
				categories[index].Recipes = append(categories[index].Recipes, product.Id)
			}
		}

		newValue := make([]interface{}, len(categories))

		for i := range categories {
			newValue[i] = categories[i]
		}

		_, err = collection.InsertMany(ctx, newValue)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	s.Logger.Info("categories seeded successfully!")

	return nil
}

func (s *Seeder) SeedMaterials(seedEntries bool) error {

	entries := []models.MaterialEntry{
		{
			Quantity:         2000,
			PurchasePrice:    250,
			PurchaseQuantity: 200,
			Company:          "Test1",
		},
		{
			Quantity:         2000,
			PurchasePrice:    250,
			PurchaseQuantity: 200,
			Company:          "Test2",
		},
		{
			Quantity:         2000,
			PurchasePrice:    250,
			PurchaseQuantity: 200,
			Company:          "Test3",
		},
		{
			Quantity:         2000,
			PurchasePrice:    250,
			PurchaseQuantity: 200,
			Company:          "Test4",
		},
		{
			Quantity:         2000,
			PurchasePrice:    250,
			PurchaseQuantity: 200,
			Company:          "Test5",
		},
	}

	material := models.Material{
		Name: "MotzarillaSeeded",
		Unit: "gm",
	}

	if seedEntries {
		material.Entries = entries
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the "components" collection from the database
	collection := client.Database("waha").Collection("components")

	// Find one document in the collection
	var component models.Material
	err = collection.FindOne(ctx, bson.M{"name": "MotzarillaSeeded"}).Decode(&component)
	if err == mongo.ErrNoDocuments {
		// Insert the materials into the database
		_, err = collection.InsertOne(ctx, material)
		if err != nil {
			return err
		}
		s.Logger.Info("materials seeded successfully")
		return nil
	} else if err != nil {
		return err
	}

	confirm_reseed_materials, err := s.Prompter.Confirmation("material already exists. Do you want to proceed with seeding materials?")
	if err != nil {
		return err
	}

	if confirm_reseed_materials {
		_, err = collection.InsertOne(ctx, material)
		if err != nil {
			return err
		}
		s.Logger.Info("materials inserted successfully")
	}

	return nil
}