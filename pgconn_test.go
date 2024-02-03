package pgconn_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/QuickDrone-Backend/pgconn"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	testConfig = pgconn.ConnConfig{
		Port: "32260",
		Host: "localhost",
		User: "testuser",
		Password: "s3cr3t",
		DbName: "testdb",
		SslMode: pgconn.SslModeDisable,
	}

	testSettings = pgconn.ConnSettings{
		MaxOpenConns: 1,
		MaxIdleConns: 1,
		ConnMaxLifetime: time.Second,
		ConnMaxIdleTime: time.Second,
	}
)

func TestConnectionInOneAttempt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	conn, err := pgconn.Connect(&testConfig, &testSettings)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	err = conn.Ping()
	require.NoError(t, err)
}

func TestConnectionInManyAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Да-да, комментарий на русском

	// В чём смысл этого теста?
	// Если Вы, дорогой читатель этого кода, заглянете в docker-compose.test.yml
	// (надеюсь я к этому времени его не переименовал)
	// то можете заметить, что я насильно ограничил количество коннектов в единицу
	// (один для рута).
	// В этом тесте в двух горутинах одновременно пытаются создать два коннекта,
	// причём один сразу (и держит коннект две секунды), а второй -- через секунду.
	// Если заменить функцию ConnectWithTries на обычную, Connect, то в любом случае
	// второй коннект потерпит ошибку и всё упадёт. Но. Именно для таких же случаев
	// нам и нужна функция ConnectWithTries, не так ли? Она делает три попытки 
	// (ну а вдруг со второй не прокатит?..)
	// и первая обязательно будет неудачной (стучится в закрытую дверь), а последующие как раз должны
	// попасть в тот момент, когда первое соединение отцепится
	
	// И да, я проверил всё ручками, так что тест рабочий

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		conn, err := pgconn.ConnectWithTries(&testConfig, &testSettings, 3, time.Second)
		require.NoError(t, err)
		require.NotNil(t, conn)

		time.Sleep(2*time.Second)

		conn.Close()
	}()
	go func() {
		defer wg.Done()

		time.Sleep(time.Second)

		conn, err := pgconn.ConnectWithTries(&testConfig, &testSettings, 3, time.Second)
		require.NoError(t, err)
		require.NotNil(t, conn)
		conn.Close()
	}()

	wg.Wait()
}

func TestMigrateUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	err := pgconn.MigrateUp(testConfig, "file:///home/zhikh/repos/my/qdrone/pgconn/test_migrations")
	if !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}

	conn, err := pgconn.Connect(&testConfig, &testSettings)
	require.NoError(t, err)

	var id int
	err = sqlx.Get(conn, &id, "SELECT id FROM test_table")
	require.NoError(t, err)
	
	require.Equal(t, id, 1)
}
