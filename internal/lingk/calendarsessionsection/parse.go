package calendarsessionsection

import (
  "io"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
)

var expectHead = []string{
  "calendarSessionExternalId",
  "courseSectionExternalId",
}

func parse(record []string) *Entry {
  return &Entry{
    Id: record[0],
    CourseSectionId: record[1],
  }
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
  entries := make([]*Entry, 0, 1024)
  errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
    entries = append(entries, parse(record))
    return nil
  })
  return entries, errs, err
}
