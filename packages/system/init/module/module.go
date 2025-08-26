package initialize

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
)

type Request struct {
	Http CustomHttp `json:"http"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
	Method       string       `json:"method"`
	Path         string       `json:"path"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
}

func Invoke(in Request) (*u.Response, error) {

	resp := u.Message(false, "0x11028:Invalid request")

	path := strings.Replace(in.Http.Path, "/", "", -1)

	if path == "nemutluturkumdiyene" {
		MigrationModels()
	} else if path == "update" {
		context := &u.Context{}

		client := &http.Client{}

		req, _ := http.NewRequest("GET", os.Getenv("VALID_API_URL"), &bytes.Buffer{})
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", in.Http.CustomHeader.Authorization)

		res, err := client.Do(req)

		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			_, errRes := u.ResMessage(false, "0x11130:Missing auth token")
			return &errRes, nil
		}

		context.UserId = uuid.MustParse(res.Header.Get("userId"))
		roleId, _ := strconv.Atoi(res.Header.Get("roleId"))
		context.RoleId = roleId
		appId, _ := strconv.Atoi(res.Header.Get("appId"))
		context.AppId = appId
		context.MerchantId = uuid.MustParse(res.Header.Get("merchantId"))
		hasId, err := strconv.ParseBool(res.Header.Get("hasId"))
		context.HasId = hasId
		projectId, _ := strconv.Atoi(res.Header.Get("projectId"))
		context.ProjectId = projectId
		context.CustomData = res.Header.Get("customData")
		initCompleted, err := strconv.ParseBool(res.Header.Get("initCompleted"))
		context.InitCompleted = initCompleted

		if context.RoleId == u.Root {
			MigrationModels()
			resp = u.Message(true, "0x11034:Migration done..")
		}
	}
	return u.Respond(resp)
}

var MigrationModels = func() {

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "services." + defaultTableName
	}

	d.GetDB().Exec("CREATE SCHEMA IF NOT EXISTS application")

	d.GetDB().Debug().AutoMigrate(&Note{})

	d.GetDB().Debug().AutoMigrate(&Bucket{}, &BucketItem{})
}
