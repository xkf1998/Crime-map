package hbase_client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
	"github.com/pierrre/geohash"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
)

type HBaseClient struct {
	env    env_interface.Env
	client gohbase.Client
	table  string
}

const (
	maxPrecision         = 12
	longitudeQualifier   = "x"
	laitudeQualifier     = "y"
	timeQualifier        = "t"
	descriptionQualifier = "d"
	yyyyMMddTHHmmss      = "2006-01-02T15:04:05"
)

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.Database.Address == "" || config.Database.Table == "" {
		return errors.New("HBase client not configured")
	}
	c, err := NewHBaseClient(env, config.Database.Address, config.Database.Table)
	if err != nil {
		return err
	}
	env.SetDatabaseClient(c)
	return nil
}

func NewHBaseClient(env env_interface.Env, zkquorum, table string) (*HBaseClient, error) {
	fmt.Printf("table name `%v`", table)
	client := gohbase.NewClient(zkquorum)
	return &HBaseClient{
		env:    env,
		client: client,
		table:  table,
	}, nil
}

func (c *HBaseClient) Conn(ctx context.Context) error {
	/*
		fmt.Println("--------- Scan test begin ----------")
		pFilter := filter.NewPrefixFilter([]byte(""))
		rangeColTFilter := filter.NewColumnRangeFilter([]byte("1970-01-19T04:12:06"), []byte("1970-01-20T02:26:25"), true, true)
		_ = rangeColTFilter
		scanRequest, err := hrpc.NewScanStr(ctx, c.table,
			hrpc.Filters(pFilter))
		scanRsp := c.client.Scan(scanRequest)
		var result *hrpc.Result
		for {
			result, err = scanRsp.Next()
			if err != nil {
				fmt.Printf("Scan error: %v\n", err)
				break
			}
			fmt.Printf("%v\n", result)
		}
		fmt.Println("--------- Scan test end ----------")
		/*
			// TODO: move the table and rowKey into config
			getRequest, err := hrpc.NewGetStr(ctx, "group4:test", "dr5rugb9rwjj1970-01-20T05:54:52")
			if err != nil {
				return fmt.Errorf("hrpc error: %v\n", err)
			}
			getRsp, err := c.client.Get(getRequest)
			if err != nil {
				return fmt.Errorf("get HBase response failed: %v\n", err)
			}
			_ = getRsp
			//fmt.Printf("DEBUG: HBase response: %v\n", getRsp)
	*/
	// c.GetCrimes(ctx, -122.3592, -122.359, 47.5272, 47.5274, 1513799100, 1593799300)
	return nil
}

func (c *HBaseClient) GetCrimes(ctx context.Context, minLongitude, maxLongitude, minLaitude, maxLaitude float64, minTime, maxTime int64) ([]*interfaces.Crime, error) {
	fmt.Printf("DEBUG: GetCrimes received minLong=%v, maxLong=%v, minLa=%v, maxLa=%v, minT=%v, maxT=%v\n",
		minLongitude, maxLongitude, minLaitude, maxLaitude, minTime, maxTime)
	if minLongitude > maxLongitude || minLaitude > maxLaitude || minTime > maxTime {
		return nil, fmt.Errorf("HBase client warnning: bad query arguments for GetCrimes")
	}

	crimes := make([]*interfaces.Crime, 0)

	minHash := geohash.Encode(minLaitude, minLongitude, maxPrecision)
	maxHash := geohash.Encode(maxLaitude, maxLongitude, maxPrecision)
	prefixRowKey := longestCommonPrefix(minHash, maxHash)
	prefixRowKeyFilter := filter.NewPrefixFilter([]byte(prefixRowKey))

	/*
		minNormalizedX := normalizeCoordinate(minLongitude, -180.0)
		maxNormalizedX := normalizeCoordinate(maxLongitude, -180.0)

		minNormalizedY := normalizeCoordinate(minLaitude, -90.0)
		maxNormalizedY := normalizeCoordinate(maxLaitude, -90.0)

		minNormalizedT := normalizeTime(minTime)
		maxNormalizedT := normalizeTime(maxTime)
	*/

	scanRequest, err := hrpc.NewScanStr(ctx, c.table, hrpc.Filters(prefixRowKeyFilter))
	scanRsp := c.client.Scan(scanRequest)

	rowCountHBaseReturned := 0
	rowCountCorrect := 0
	var result *hrpc.Result
	for {
		result, err = scanRsp.Next()
		if err != nil {
			break
		}
		crime := &interfaces.Crime{}
		for _, cell := range result.Cells {
			switch string(cell.Qualifier) {
			case longitudeQualifier:
				crime.Longitude = denormalizeCoordinate(string(cell.Value), -180.0)
			case laitudeQualifier:
				crime.Latitude = denormalizeCoordinate(string(cell.Value), -90.0)
			case timeQualifier:
				crime.Time = denormalizeTime(string(cell.Value))
			case descriptionQualifier:
				crime.Description = string(cell.Value)
			default:
				// unexpeced qualifier: just skip it to fit changes in schema
			}
		}
		rowCountHBaseReturned++
		// TODO: remove these condition stmts after filter implemented
		if crime.Longitude < minLongitude || crime.Longitude > maxLongitude || crime.Latitude < minLaitude || crime.Latitude > maxLaitude {
			continue
		}
		rowCountCorrect++
		if crime.Time < minTime || crime.Time > maxTime {
			continue
		}
		// fmt.Printf("%v\n", *crime)
		crimes = append(crimes, crime)
	}
	fmt.Printf("DEBUG: query-prefix length %v out of 12, prefix-match hit rate %.2f%% out of %v total returned records\n",
		len(prefixRowKey), 100*float64(rowCountCorrect)/float64(rowCountHBaseReturned), len(crimes))
	return crimes, nil
}

func (c *HBaseClient) Close() error {
	c.client.Close()
	return nil
}

func longestCommonPrefix(s1, s2 string) string {
	for i := 0; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			return s1[:i]
		}
	}
	if len(s1) < len(s2) {
		return s1
	} else {
		return s2
	}
}

func normalizeCoordinate(value, minValue float64) string {
	posValue := value - minValue
	prefixValue := int64(posValue)
	prefixStr := strconv.FormatInt(prefixValue, 10)
	prefixPadding := strings.Repeat("0", 3-len(prefixStr))
	prefixStr = prefixPadding + prefixStr

	suffixValue := int64((posValue - float64(prefixValue)) * 1e6)
	suffixStr := strconv.FormatInt(suffixValue, 10)
	suffixPadding := strings.Repeat("0", 6-len(suffixStr))
	suffixStr += suffixPadding

	return prefixStr + suffixStr
}

func denormalizeCoordinate(normalizedStr string, minValue float64) float64 {
	// TODO: Don't ignore this error
	value, _ := strconv.ParseFloat(normalizedStr, 64)
	return value/1e6 + minValue
}

func normalizeTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(yyyyMMddTHHmmss)
}

func denormalizeTime(normalizeTime string) int64 {
	// TODO: Don't ignore this error
	date, _ := time.Parse(yyyyMMddTHHmmss, normalizeTime)
	return date.Unix()
}
