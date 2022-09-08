package main

import (
	"fmt"
	"os"

	"github.com/teal-finance/quid/server/db"
)

func writeQuidAdminToken(username string) error {
	return writeDevAdminToken(username, "quid")
}

func writeNsAdminToken(username, namespace string) error {
	return writeDevAdminToken(username, namespace)
}

func writeDevAdminToken(username, namespace string) error {
	// generate a refresh token
	token, err := db.GenNsAdminTokenForUser(username, namespace)
	if err != nil {
		msg := "Error generating refresh token"
		log.QueryError(msg, err)
		return err
	}

	// write the ui env file

	dir, err := os.Getwd()
	if err != nil {
		log.Error(err)
		return err
	}
	relpath := "/ui/.env.development.local"
	filepath := dir + relpath
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	v := "VITE_DEV_TOKEN=\"" + token + "\""
	_, err = f.Write([]byte(v))
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	fmt.Println("Dev token written in", relpath)
	return nil
}
