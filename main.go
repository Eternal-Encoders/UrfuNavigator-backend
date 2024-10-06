package main

import (
	"UrfuNavigator-backend/internal/app"
	database "UrfuNavigator-backend/internal/database/mongo"
	"UrfuNavigator-backend/internal/objstore"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file. Env load from system only")
	}
}

func main() {
	uri, exist := os.LookupEnv("DATABASE_URI")
	if !exist {
		file, exist := os.LookupEnv("DATABASE_URI_FILE")

		if !exist {
			log.Fatal("No connection uri")
		}

		data, err := os.ReadFile(file)

		if err != nil {
			log.Fatal(err)
		}

		uri = string(data)
	}

	collection, exist := os.LookupEnv("DATABASE_COLLECTION")
	if !exist {
		log.Fatal("No collection specified")
	}

	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("No port specified")
	}

	cors, exist := os.LookupEnv("CORS")
	if !exist {
		log.Fatal("Cors policy not specified")
	}

	s3Endpoint, exist := os.LookupEnv("BUCKET_ENDPOINT")
	if !exist {
		log.Fatal("No s3 endpoint specified")
	}

	s3Access, exist := os.LookupEnv("BUCKET_ACCESS_KEY")
	if !exist {
		file, exist := os.LookupEnv("BUCKET_ACCESS_KEY_FILE")

		if !exist {
			log.Fatal("No s3 access key specified")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		s3Access = string(data)
	}

	s3Secret, exist := os.LookupEnv("BUCKET_SECRET_KEY")
	if !exist {
		file, exist := os.LookupEnv("BUCKET_SECRET_KEY_FILE")

		if !exist {
			log.Fatal("No s3 secret key specified")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		s3Secret = string(data)
	}

	bucketName, exist := os.LookupEnv("BUCKET_NAME")
	if !exist {
		log.Fatal("No s3 bucket specified")
	}

	db := database.Connect(uri, collection)
	objectStore := objstore.Connect(s3Endpoint, s3Access, s3Secret, bucketName)
	api := app.NewAPI(
		port,
		db,
		objectStore,
		cors,
	)

	defer db.Disconnect()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
