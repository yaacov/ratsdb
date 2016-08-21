package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yaacov/ratsdb/models"
)

// Open DB create tables and start transaction
func init() {
	var db *sql.DB
	var err error

	// open db
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create data table
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// open transaction
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
}

// build sql for select query
func SelectSqlQuery(key string, start string, end string, labels string) string {
	var where bytes.Buffer
	var whereLabels bytes.Buffer

	// create a where
	where.WriteString(fmt.Sprintf(" where id is not null"))
	startTime, err := strconv.ParseUint(start, 0, 64)
	if err == nil {
		where.WriteString(fmt.Sprintf(" and time >= %d", startTime))
	}
	endTime, err := strconv.ParseUint(end, 0, 64)
	if err == nil {
		where.WriteString(fmt.Sprintf(" and time < %d", endTime))
	}
	keyOk, err := regexp.MatchString("^[a-zA-Z0-9-]+$", key)
	if keyOk && err == nil {
		where.WriteString(fmt.Sprintf(" and key = \"%s\"", key))
	}
	labelsList := strings.Split(labels, ",")
	for _, label := range labelsList {
		labelOk, err := regexp.MatchString("^[a-zA-Z0-9-]+$", label)
		if labelOk && err == nil {
			whereLabels.WriteString(fmt.Sprintf(" and labels like \"%%%s%%\"", label))
		}
	}

	return where.String() + whereLabels.String()
}

// Get list of data objects from table
func Data(key string, start string, end string, labels string) models.Samples {
	var samples models.Samples

	// prepare query statement
	stmt, err := tx.Prepare(selectFromTable + SelectSqlQuery(key, start, end, labels))
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// query samples
	rows, err := stmt.Query()
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	// append samples to output list
	for rows.Next() {
		var sample models.Sample
		err = rows.Scan(&sample.Id, &sample.Time, &sample.Key, &sample.Value, &sample.Labels)
		if err != nil {
			log.Print(err)
		}
		samples = append(samples, sample)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
	}

	return samples
}

// Get list of data objects from table
func DataBuckets(key string, start string, end string, labels string, bucket string) models.Buckets {
	var buckets models.Buckets

	bucketTime, err := strconv.ParseUint(bucket, 0, 64)
	if err != nil {
		// default bucket time is 100 * Millisecond
		bucketTime = 100 * 1000
	}

	// prepare query statement
	stmt, err := tx.Prepare(selectBucketFromTable + SelectSqlQuery(key, start, end, labels) + " group by key, (time / ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// query samples
	rows, err := stmt.Query(bucketTime)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	// append samples to output list
	for rows.Next() {
		var bucket models.Bucket
		err = rows.Scan(&bucket.Key, &bucket.Count, &bucket.Start, &bucket.End, &bucket.Min, &bucket.Max, &bucket.Avg)
		if err != nil {
			log.Print(err)
		}
		buckets = append(buckets, bucket)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
	}

	return buckets
}

// Find one data object in table
func DataFind(id int) models.Sample {
	var sample models.Sample

	// prepare query statement
	stmt, err := tx.Prepare(selectOneFromTable)
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// query one sample
	err = stmt.QueryRow(id).Scan(&sample.Id, &sample.Time, &sample.Key, &sample.Value, &sample.Labels)
	if err != nil {
		log.Print(err)
	}

	return sample
}

// Insert a data object to table
func DataCreate(t models.Sample) models.Sample {
	// get currentId
	currentId += 1

	// set the new sample Id and Time fields
	t.Id = currentId
	t.Time = time.Now().UnixNano() / int64(time.Millisecond)

	// prepare insert statement
	stmt, err := tx.Prepare(insertToTable)
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// insert data object to table
	_, err = stmt.Exec(t.Id, t.Time, t.Key, t.Value, t.Labels)
	if err != nil {
		log.Print(err)
	}

	return t
}

// global variables
var (
	currentId int
	tx        *sql.Tx
)

const createTable = `create table if not exists samples (
  id integer not null primary key,
  time numeric,
  key text,
  value numeric,
  labels text
);`
const insertToTable = `insert into samples (id, time, key, value, labels) values (?, ?, ?, ?, ?);`
const selectOneFromTable = `select id, time, key, value, labels from samples where id = ?;`
const selectFromTable = `select id, time, key, value, labels from samples`
const selectBucketFromTable = `select key, count(id) as count,
  min(time) as start, max(time) as end,
	min(value) as min, max(value) as max, avg(value) as avg from samples`
