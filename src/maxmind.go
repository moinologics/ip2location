package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/oschwald/geoip2-golang"
)

type LocationRecord struct {
	IPType      string        `json:"ipType"`
	Location    *geoip2.City  `json:"location"`
	ElapsedTime time.Duration `json:"elapsedTime"`
}

const MAXMIND_GEOIP_CONF_DIR = "/usr/local/etc"
const MAXMIND_GEOIP_CONF_FILENAME = "GeoIP.conf"
const MAXMIND_GEOIP_DB_DIR = "/usr/local/share/GeoIP"

var editionIds = strings.Split(getEnv("MAXMIND_EDITION_IDS", ""), " ")

func setupConfigFile() {

	sampleConfFileName := fmt.Sprintf("%s/%s.sample", MAXMIND_GEOIP_CONF_DIR, MAXMIND_GEOIP_CONF_FILENAME)
	confFileName := fmt.Sprintf("%s/%s", MAXMIND_GEOIP_CONF_DIR, MAXMIND_GEOIP_CONF_FILENAME)

	sampleConfContent, err := os.ReadFile(sampleConfFileName)
	if err != nil {
		log.Fatalf("Error reading sample conf file: %v", err)
	}

	accountId := getEnv("MAXMIND_ACCOUNT_ID", "")
	licenseKey := getEnv("MAXMIND_LICENSE_KEY", "")
	editionIDs := getEnv("MAXMIND_EDITION_IDS", "")

	if accountId == "" || licenseKey == "" || editionIDs == "" {
		log.Fatalln("at least one of env MAXMIND_ACCOUNT_ID, MAXMIND_LICENSE_KEY, MAXMIND_EDITION_IDS are missing or empty")
	}

	confFileContent := strings.Replace(string(sampleConfContent), "<account-id>", accountId, 1)
	confFileContent = strings.Replace(string(confFileContent), "<licenseKey-key>", licenseKey, 1)
	confFileContent = strings.Replace(string(confFileContent), "<addition-ids>", editionIDs, 1)

	err = os.WriteFile(confFileName, []byte(confFileContent), 0644)
	if err != nil {
		log.Fatalf("Error writing conf file: %v", err)
	}

	log.Printf("conf file %s created from %s\n", confFileName, sampleConfFileName)
}

func updateDB() error {
	stderr := bytes.Buffer{}
	cmd := exec.Command("geoipupdate")
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Error When Running database update: %s%v", string(stderr.Bytes()), err)
		return err
	}
	log.Println("database update successfull")
	return nil
}

func getIpLocation(inputIP string, dbEditionId string) (*LocationRecord, error) {

	start := time.Now()

	db, err := geoip2.Open(fmt.Sprintf("%s/%s.mmdb", MAXMIND_GEOIP_DB_DIR, dbEditionId))

	if err != nil {
		log.Printf("error while opening database %s: %v", dbEditionId, err)
		return nil, err
	}
	defer db.Close()

	ip := net.ParseIP(inputIP)

	if ip == nil {
		log.Printf("invalid IP %v provided", ip)
		return nil, errors.New("invalid IP provided")
	}

	if ip.IsPrivate() {
		log.Printf("private IP %v provided", ip)
		return nil, errors.New("private IPs is not allowed")
	}

	record, err := db.City(ip)
	if err != nil {
		log.Printf("error when searching ip in db: %v", err)
		return nil, errors.New("error when searching ip")
	}

	if len(record.Country.Names) == 0 {
		log.Printf("no record found for ip %v\n", ip)
		return nil, errors.New("not found")
	}

	ipType := "IPv4"

	if ip.To4() == nil {
		ipType = "IPv6"
	}

	elapsed := time.Since(start)

	result := LocationRecord{IPType: ipType, Location: record, ElapsedTime: elapsed}

	// resultBytes, err := json.Marshal(result)
	// if err != nil {
	// 	log.Printf("error when parsing record from db %v", err)
	// 	return nil, errors.New("error when parsing record from db")
	// }

	// _ = string(resultBytes)

	return &result, nil
}

func ip2location(ctx echo.Context) error {
	apiKey := ctx.Request().Header.Get("API-KEY")

	allowedAPIKey := getEnv("ALLOWED_API_KEY", "")

	if allowedAPIKey != "" && apiKey != allowedAPIKey {
		ctx.String(http.StatusUnauthorized, "invalid api key")
	}

	dbEdition := ctx.Request().Header.Get("MAXMIND-DB-EDITION")

	if dbEdition != "" {
		if !slices.Contains(editionIds, dbEdition) {
			ctx.String(http.StatusBadRequest, "invalid maxmind db edition")
		}
	} else {
		dbEdition = editionIds[0]
	}

	ip := ctx.Param("ip")

	result, err := getIpLocation(ip, dbEdition)

	if err == nil {
		return ctx.JSON(http.StatusOK, result)
	}

	errMsg := err.Error()

	if strings.HasPrefix(errMsg, "error") {
		return ctx.String(http.StatusInternalServerError, "something went wrong")
	}

	if errMsg == "not found" {
		return ctx.String(http.StatusNotFound, errMsg)
	}

	return ctx.String(http.StatusBadRequest, errMsg)
}
