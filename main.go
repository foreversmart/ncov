package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type City struct {
	Province string  `json:"province"`
	City     string  `json:"city"`
	Data     []*Data `json:"data"`
}

type Province struct {
	Province string  `json:"province"`
	Data     []*Data `json:"data"`
}

type Data struct {
	Confirmed   int       `json:"confirmed"`
	Suspected   int       `json:"suspected"`
	Cured       int       `json:"cured"`
	Dead        int       `json:"dead"`
	UpdatedTime time.Time `json:"updated_time"`
}

const timeFormat = "2006-01-02 15:04:05.000"

func main() {
	pwd, _ := os.Getwd()
	filename := filepath.Join(pwd, "DXYArea.csv")

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		println("open file", filename, err)
		return
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		println("read csv data", err)
		return
	}

	cityMap := make(map[string]*City)
	provinceMap := make(map[string]*Province)
	for i, line := range records {
		if i == 0 {
			continue
		}

		var (
			provinceName, cityName                                                                string
			provinceConfirmedCount, provinceSuspectedCount, provinceCuredCount, provinceDeadCount string
			cityConfirmedCount, citySuspectedCount, cityCuredCount, cityDeadCount                 string
			updateTime                                                                            string
		)

		var (
			provinceConfirmed, provinceSuspected, provinceCured, provinceDead int64
			cityConfirmed, citySuspected, cityCured, cityDead                 int64
		)
		for j, item := range line {
			switch j {
			case 0:
				provinceName = item
			case 1:
				cityName = item
			case 2:
				provinceConfirmedCount = item
				provinceConfirmed, _ = strconv.ParseInt(provinceConfirmedCount, 10, 64)
			case 3:
				provinceSuspectedCount = item
				provinceSuspected, _ = strconv.ParseInt(provinceSuspectedCount, 10, 64)
			case 4:
				provinceCuredCount = item
				provinceCured, _ = strconv.ParseInt(provinceCuredCount, 10, 64)
			case 5:
				provinceDeadCount = item
				provinceDead, _ = strconv.ParseInt(provinceDeadCount, 10, 64)
			case 6:
				cityConfirmedCount = item
				cityConfirmed, _ = strconv.ParseInt(cityConfirmedCount, 10, 64)
			case 7:
				citySuspectedCount = item
				citySuspected, _ = strconv.ParseInt(citySuspectedCount, 10, 64)
			case 8:
				cityCuredCount = item
				cityCured, _ = strconv.ParseInt(cityCuredCount, 10, 64)
			case 9:
				cityDeadCount = item
				cityDead, _ = strconv.ParseInt(cityDeadCount, 10, 64)
			case 10:
				updateTime = item

			}
		}

		updated, _ := time.ParseInLocation(timeFormat, updateTime, time.FixedZone("UTC", 8*60*60))

		cityName = GetCityName(cityName, provinceName)
		provinceData := &Data{
			Confirmed:   int(provinceConfirmed),
			Suspected:   int(provinceSuspected),
			Dead:        int(provinceDead),
			Cured:       int(provinceCured),
			UpdatedTime: updated,
		}

		var cityData *Data
		if strings.HasSuffix(provinceName, "市") {
			cityName = provinceName
			cityData = provinceData

		} else {
			cityData = &Data{
				Confirmed:   int(cityConfirmed),
				Suspected:   int(citySuspected),
				Dead:        int(cityDead),
				Cured:       int(cityCured),
				UpdatedTime: updated,
			}

		}

		if _, ok := cityMap[cityName]; !ok {
			cityMap[cityName] = &City{
				Province: provinceName,
				City:     cityName,
				Data:     make([]*Data, 0, 10),
			}
		}

		cityMap[cityName].Data = append(cityMap[cityName].Data, cityData)

		if _, ok := provinceMap[provinceName]; !ok {
			provinceMap[provinceName] = &Province{
				Province: provinceName,
				Data:     make([]*Data, 0, 10),
			}
		}

		provinceMap[provinceName].Data = append(provinceMap[provinceName].Data, provinceData)
	}

	buf := &bytes.Buffer{}
	resarr := make([]int, 0, 10)
	for city, v := range cityMap {
		res := v.Calc()
		resarr = append(resarr, res)

		buf.WriteString(fmt.Sprintf("%s&%d&", city, v.Calc()))
	}

	sort.Ints(resarr)
	l := len(resarr)
	fmt.Println("res city length", l)
	fmt.Println("city score rank：")
	fmt.Println(resarr[int(float32(l)*0.3)])
	fmt.Println(resarr[int(float32(l)*0.5)])
	fmt.Println(resarr[int(float32(l)*0.65)])
	fmt.Println(resarr[int(float32(l)*0.75)])
	fmt.Println(resarr[int(float32(l)*0.80)])
	fmt.Println(resarr[int(float32(l)*0.85)])
	fmt.Println(resarr[int(float32(l)*0.90)])
	fmt.Println(resarr[int(float32(l)*0.94)])
	fmt.Println(resarr[int(float32(l)*0.98)])
	fmt.Println(resarr[int(float32(l)*0.99)])

	err = ioutil.WriteFile(filepath.Join(pwd, "data.txt"), buf.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

}

const score = 5

func (c *City) Calc() int {
	day := -1
	arr := make([]*Data, 0, 5)
	for _, d := range c.Data {
		if day == d.UpdatedTime.Day() {
			continue
		}

		day = d.UpdatedTime.Day()

		arr = append(arr, d)
	}

	count := 0
	for i := len(arr) - 1; i >= 0; i-- {
		if i == len(arr)-1 {
			count = arr[i].Confirmed * score
			continue
		}

		delta := arr[i].Confirmed - arr[i+1].Confirmed
		count = count + delta*score
		for j := i + 1; j < len(arr)-1; j++ {
			if j-i > 10 {
				break
			}

			count = count - (arr[j+1].Confirmed - arr[j].Confirmed)

			if j == len(arr)-2 && j+1-i > 10 {
				count = count - arr[j+1].Confirmed
			}
		}

		// check bonus
		if i < len(arr)-2 {
			d1 := arr[i].Confirmed - arr[i+1].Confirmed
			d2 := arr[i+1].Confirmed - arr[i+2].Confirmed

			count = count - d2 + d1
			if d1 == d2 {
				count = count - 1
			}
		}

	}

	return count

}

func GetCityName(cityName, provinceName string) (res string) {
	switch cityName {
	case "神农架林区":
		return cityName
	case "恩施", "恩施州":
		return "恩施土家族苗族自治州"
	case "湘西自治州":
		return "湘西土家族苗族自治州"
	case "大兴安岭":
		return "大兴安岭地区"
	case "黔东南州":
		return "黔东南苗族侗族自治州"
	case "黔西南州":
		return "黔西南布依族苗族自治州"
	case "黔南州":
		return "黔南布依族苗族自治州"
	case "兴安盟乌兰浩特":
		return "兴安盟"
	case "阿坝州":
		return "阿坝藏族羌族自治州"
	case "甘孜州":
		return "甘孜藏族自治州"
	case "凉山州":
		return "凉山彝族自治州"
	case "西双版纳":
		return "西双版纳傣族自治州"
	case "德宏":
		return "德宏傣族景颇族自治州"
	case "大理":
		return "大理白族自治州"
	case "红河":
		return "红河哈尼族彝族自治州"
	case "伊犁州":
		return "伊犁哈萨克自治州"
	case "阿克苏":
		return "阿克苏地区"
	case "文山":
		return "文山壮族苗族自治州"
	case "楚雄":
		return "楚雄彝族自治州"
	case "琼中县":
		return "琼中黎族苗族自治县"
	case "定安":
		return "定安县"
	case "陵水县":
		return "陵水黎族自治县"
	case "昌江":
		return "昌江黎族自治县"
	case "乐东":
		return "乐东黎族自治县"
	case "临夏":
		return "临夏回族自治州"
	case "甘南":
		return "甘南藏族自治州"
	}

	if strings.HasSuffix(cityName, "县") {
		return cityName
	}

	if strings.Contains(cityName, "锡林郭勒盟") {
		return "锡林郭勒盟"
	}

	if !strings.HasSuffix(provinceName, "市") && !strings.HasSuffix(cityName, "自治州") && !strings.HasSuffix(cityName, "市") {
		res = cityName + "市"
	} else {
		res = cityName
	}

	if strings.Contains(res, "(") {
		items := strings.Split(res, "(")
		aItems := strings.Split(items[1], ")")
		res = items[0] + aItems[1]
		return
	}

	if strings.Contains(res, "（") {
		items := strings.Split(res, "（")
		aItems := strings.Split(items[1], "）")
		res = items[0] + aItems[1]
		return
	}

	return
}
