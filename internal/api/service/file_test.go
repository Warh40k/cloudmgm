package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestFileService_UploadFile_HappyPath(t *testing.T) {
	/**
	Создать таблицу кейсов
	TODO
	*/

	//load config
	pathToRoot := "../../../"
	viper.AddConfigPath(pathToRoot + "configs")
	viper.SetConfigName("local")
	err := viper.ReadInConfig()
	require.NoError(t, err)

	volumeId := uuid.New()
	fileName := gofakeit.Sentence(5) + gofakeit.Extension()
	fs := afero.NewMemMapFs()
	file, err := fs.Create(fileName)
	require.NoError(t, err)
	volumePath := viper.GetString("files.save_path") + "/" + volumeId.String()

	//init layers
	db, _, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := repository.NewRepository(dbx, log)
	service := NewService(repo, log)

	// run upload
	saveName, err := service.UploadFile(volumePath, file, fileName, fs)
	assert.Equal(t, fileName, saveName)

	ok, err := afero.Exists(fs, viper.GetString("files.save_path")+
		"/"+volumeId.String()+"/"+fileName)
	assert.True(t, ok)
}
