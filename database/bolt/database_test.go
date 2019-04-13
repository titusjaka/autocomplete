package bolt

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/logutils"
	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/autocomplete"
	"github.com/titusjaka/autocomplete/database"
)

func TestDatabaseReadWrite(t *testing.T) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer:   os.Stdout,
	}
	logger := log.New(filter, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	boltDb, err := New("test.db")
	if err != nil {
		logger.Fatalf("[FATAL] failed to open bolt db: %v", err)
	}
	assert.NotNil(t, boltDb)
	defer os.Remove("test.db")

	actualWD, err := boltDb.Get("test-request")
	assert.Equal(t, database.ErrNotFound, err)
	assert.Empty(t, actualWD)

	testWD := []*autocomplete.WidgetData{
		{
			Slug:     "TST",
			Title:    "Test",
			Subtitle: "Yep, it's really just a test",
		},
		{
			Slug:     "LOL",
			Title:    "Loling",
			Subtitle: "Test data 2",
		},
	}

	err = boltDb.Write("test-request", testWD)
	assert.NoError(t, err)

	actualWD, err = boltDb.Get("test-request")
	assert.NoError(t, err)
	assert.Equal(t, testWD, actualWD)
}
