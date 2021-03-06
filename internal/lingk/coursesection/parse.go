package coursesection

import (
	"errors"
	"github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"io"
	"math"
	"strconv"
)

var expectHead = []string{
	"externalId",
	"courseExternalId",
	"coureSectionNumber", // [sic]
	"capacity",
	"currentEnrollment",
	"Status",
	"CreditHours",
}

func parse(record []string) (*Entry, error) {
	colExternalId := record[0]
	colCourseExternalId := record[1]
	colCourseSectionNumber := record[2]
	colCapacity := record[3]
	colCurrentEnrollment := record[4]
	colStatus := record[5]
	colCreditHours := record[6]

	section, err := strconv.Atoi(colCourseSectionNumber)
	if err != nil {
		return nil, errors.New("invalid section")
	}

	seats, err := data.ParseSeats(colCurrentEnrollment, colCapacity)
	if err != nil {
		return nil, errors.New("invalid seats")
	}

	status, err := data.ParseStatus(colStatus)
	if err != nil {
		return nil, errors.New("invalid status")
	}

	credits, err := strconv.ParseFloat(colCreditHours, 32)
	if err != nil {
		return nil, errors.New("invalid float credits")
	}
	quarterCredits := credits * 4
	if math.Round(quarterCredits) != quarterCredits {
		// We shouldn't have to worry about floating-point precision issues here,
		// since all "valid" credits should be some multiple of 0.25, which is
		// dyadic and therefore can be represented exactly in binary.  So it
		// suffices to directly check floating-point equality; if any issues arise,
		// the credit count must not actually have been a multiple of 0.25 in the
		// first place and would be unexpected/"invalid" anyway.
		return nil, errors.New("float credits unexpected value")
	}

	return &Entry{
		Id:             colExternalId,
		CourseId:       colCourseExternalId,
		Section:        section,
		Seats:          seats,
		Status:         status,
		QuarterCredits: int(math.Round(quarterCredits)),
	}, nil
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
	courseSections := make([]*Entry, 0, 1024)
	errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
		courseSection, err := parse(record)
		if err != nil {
			return err
		}
		courseSections = append(courseSections, courseSection)
		return nil
	})
	return courseSections, errs, err
}
