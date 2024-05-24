package db

import (
	"fmt"
	"log"
)

type DBClient struct {
	Prisma *PrismaClient
}

func NewPrismaClient() *DBClient {
	prismaClient := NewClient()
	return &DBClient{
		Prisma: prismaClient,
	}
}

func (db *DBClient) Connect() error {
	err := db.Prisma.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Connected to database")
	return nil
}

func (db *DBClient) Disconnect() error {
	err := db.Prisma.Disconnect()
	if err != nil {
		return fmt.Errorf("failed to disconnect from database: %w", err)
	}
	log.Println("Disconnected from database")
	return nil
}
