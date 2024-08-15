package query

import (
	"fmt"
	"math"
	"strconv"
	"time"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductFilters(mode string, req *pb.GetProductFilter) (filter, pagination, sorting string, earliest, latest *timestamppb.Timestamp, err error) {
	if req.VisEarly == "" {
		earliest = timestamppb.New(time.Unix(0, 0))
	} else {
		parsedTime, err := time.Parse(time.DateOnly, req.VisEarly)
		if err != nil {
			return "", "", "", nil, nil, fmt.Errorf("failed to parse VisEarly: %v", err)
		}
		earliest = timestamppb.New(parsedTime)
	}
	if req.VisLate == "" {
		latest = timestamppb.New(time.Unix(1<<63-62135596801, 999999999))
	} else {
		parsedTime, err := time.Parse(time.DateOnly, req.VisLate)
		if err != nil {
			return "", "", "", nil, nil, fmt.Errorf("failed to parse VisLate: %v", err)
		}
		latest = timestamppb.New(parsedTime)
	}

	if req.Visibility != "" {
		filter += " AND visibility = '" + req.Visibility + "'"
	}

	if req.Exclusion != "" {
		filter += " AND exclusion = '" + req.Exclusion + "'"
	}

	if req.IsAdminVerified != "" {
		filter += " AND is_admin_verified = '" + req.IsAdminVerified + "'"
	}

	if req.HighestPrice == 0 {
		req.HighestPrice = math.MaxInt32
	} else {
		filter += " AND price BETWEEN " + strconv.Itoa(int(req.LowestPrice)) + " AND " + strconv.Itoa(int(req.HighestPrice))
	}

	if len(req.Categories) > 0 {
		filter += " AND category_name IN ("
		for _, category := range req.Categories {
			filter = "'" + category.Name + "', "
		}
		filter = filter[:len(filter)-2] + ")"
	}

	if req.Pagination == (&pb.Pagination{}) {
		req.Pagination = &pb.Pagination{Page: 1, Limit: 10}
	}
	req.Pagination.Offset = (req.Pagination.Page - 1) * req.Pagination.Limit
	pagination = "LIMIT " + strconv.Itoa(int(req.Pagination.Limit)) + " OFFSET " + strconv.Itoa(int(req.Pagination.Offset))

	if req.Sorting.OrderBy != "" {
		sorting = " ORDER BY " + req.Sorting.OrderBy
		if req.Sorting.IsReversed {
			sorting += " DESC"
		} else {
			sorting += " ASC"
		}
	}

	if req.Search != "" && (mode == "user" || mode == "admin") {
		filter = " AND name ILIKE '%" + req.Search + "%' OR partner_name ILIKE '%" + req.Search + "%'"
	} else if req.Search != "" && mode == "seller" {
		filter = " AND name ILIKE '%" + req.Search
	}

	return filter, pagination, sorting, earliest, latest, nil
}
