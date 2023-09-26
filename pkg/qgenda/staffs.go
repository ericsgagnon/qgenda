package qgenda

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"

	"github.com/jmoiron/sqlx"
)

type Staffs []Staff

func GetStaffs(ctx context.Context, c *Client, rc *RequestConfig) (Staffs, error) {
	req := NewStaffRequest(rc)
	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ss := Staffs{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return nil, err
	}
	sourceQuery := resp.Request.URL.String()
	extractDateTime, err := ParseTime(resp.Header.Get("date"))
	if err != nil {
		return nil, err
	}

	if len(ss) > 0 {
		for i, _ := range ss {
			ss[i].SourceQuery = &sourceQuery
			ss[i].ExtractDateTime = &extractDateTime
		}
	}
	log.Printf("staffs:\ttotal: %d\textractDateTime:%s\n", len(ss), extractDateTime)
	return ss, nil
}

func GetStaffsFromFiles(filenames ...string) (Staffs, error) {
	staffs := Staffs{}
	for _, filename := range filenames {

		fi, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}
		modTime := fi.ModTime()

		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &staffs); err != nil {
			return nil, err
		}
		for i, v := range staffs {
			if v.ExtractDateTime == nil {
				proxyExtractDateTime := NewTime(modTime)
				v.ExtractDateTime = &proxyExtractDateTime
			}
			staffs[i] = v
		}
	}
	return staffs, nil
}

func (s *Staffs) Get(ctx context.Context, c *Client, rc *RequestConfig) error {
	switch staffs, err := GetStaffs(ctx, c, rc); {
	case err != nil:
		return err
	default:
		*s = staffs
		return nil
	}
}

func (s *Staffs) Process() error {
	ss := *s
	ss.Sort()
	for i, _ := range ss {
		if err := ss[i].Process(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Staffs) Sort() *Staffs {
	ss := *s
	sort.SliceStable(ss, func(i, j int) bool {
		itime := *(ss[i].ExtractDateTime)
		jtime := *(ss[j].ExtractDateTime)
		ikey := *(ss[i].StaffKey)
		jkey := *(ss[j].StaffKey)
		switch {
		case itime.Time.Before(jtime.Time):
			return true
		case !itime.Time.After(jtime.Time) && ikey < jkey:
			return true
		default:
			return false
		}
	})
	*s = ss
	return s
}

func (s Staffs) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return CreatePGTable(ctx, db, s, schema, table)
}

func (s Staffs) PutPG(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return PutPG(ctx, db, s, schema, table)
}

func (s Staffs) GetPGStatus(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	return GetPGStatus(ctx, db, s, schema, table, "")
}
