package api

import (
	"encoding/json"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strings"
)

func (me *T) GetNNSNameByAdmin(args struct {
	Asset  h160.T
	Admin  h160.T
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Admin.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	adminstr := string(args.Admin)
	little_endian := helper.HexToBytes(adminstr[2:len(adminstr)])

	rea := helper.ReverseBytes(little_endian)
	encodeAdmin := crypto.Base64Encode(rea)
	bakEncodeAdmin := ""
	if strings.Index(encodeAdmin, "+") >= 0 {
		bakEncodeAdmin = strings.Split(encodeAdmin, "+")[0] + "+"
	}

	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Nep11Properties",
			Index:      "GetNFTByWords",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": args.Asset}},
				bson.M{"$match": bson.M{"$or": []interface{}{
					bson.M{"properties": bson.M{"$regex": "admin\":\"" + encodeAdmin, "$options": "$i"}},
					bson.M{"properties": bson.M{"$regex": "admin\": \"" + encodeAdmin, "$options": "$i"}},
					bson.M{"properties": bson.M{"$regex": "admin\":\"" + bakEncodeAdmin, "$options": "$i"}},
					bson.M{"properties": bson.M{"$regex": "admin\": \"" + bakEncodeAdmin, "$options": "$i"}},
				}}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}
	var namelist []string
	for _, item := range r1 {
		//获取nft 属性
		if item["properties"] != nil {
			extendData := item["properties"].(string)
			if extendData != "" {
				var data map[string]interface{}
				if err2 := json.Unmarshal([]byte(extendData), &data); err2 == nil {
					name, ok := data["name"]
					if ok {
						item["name"] = name
						namelist = append(namelist, name.(string))
					}
					admin, ok2 := data["admin"]
					if ok2 {
						item["admin"] = admin
					}
					expiration, ok3 := data["expiration"]
					if ok3 {
						item["expiration"] = expiration
					}

				} else {
					return err2
				}

			}
		}

	}

	result := make([]map[string]interface{}, 0)
	result = append(result, map[string]interface{}{"name": namelist})
	//r3, err := me.FilterAggragateAndAppendCount(result, count, args.Filter)

	//
	if err != nil {
		return err
	}
	r, err := json.Marshal(result)
	if err != nil {
		return err
	}

	*ret = json.RawMessage(r)
	return nil
}

func UnicodeRegexToString(source string) string {

	return source
}
