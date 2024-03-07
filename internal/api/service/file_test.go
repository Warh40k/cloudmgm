package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

func LoadConfig(t *testing.T) {
	pathToRoot := "../../../"
	viper.AddConfigPath(pathToRoot + "configs")
	viper.SetConfigName("local")
	err := viper.ReadInConfig()
	require.NoError(t, err)
}

func InitLayers(t *testing.T) *Service {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	dbx := sqlx.NewDb(db, "sqlmock")
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := repository.NewRepository(dbx, log)
	service := NewService(repo, log)

	return service
}

func TestFileService_UploadFile(t *testing.T) {
	LoadConfig(t)
	service := InitLayers(t)

	fs := afero.NewMemMapFs()

	testTable := []struct {
		Label         string
		NumIters      int
		Path          string
		VolumeId      string
		FileName      string
		FileExtension string
		ExpectError   bool
	}{
		{
			Label:         "HappyPath",
			NumIters:      1,
			Path:          viper.GetString("files.save_path"),
			VolumeId:      gofakeit.UUID(),
			FileName:      gofakeit.Sentence(5),
			FileExtension: gofakeit.Extension(),
			ExpectError:   false,
		},
		{
			Label:       "EmptyPath",
			NumIters:    1,
			Path:        "",
			VolumeId:    "",
			FileName:    "",
			ExpectError: true,
		},
	}

	for _, tbl := range testTable {
		t.Run(tbl.Label, func(t *testing.T) {
			filename := tbl.FileName + tbl.FileExtension
			file, err := fs.Create(filename)
			require.NoError(t, err)

			volumePath := tbl.Path + "/" + tbl.VolumeId
			var saveName string
			for j := 0; j < tbl.NumIters; j++ {
				saveName, err = service.UploadFile(volumePath, file, filename, fs)
				if tbl.ExpectError {
					require.Error(t, err)
					return
				}
				namePostfix := ""
				if j != 0 {
					namePostfix = "(" + strconv.Itoa(j) + ")"
				}
				expectedName := tbl.FileName + namePostfix + tbl.FileExtension
				ok, err := afero.Exists(fs, volumePath+"/"+saveName)
				assert.NoError(t, err)
				assert.Equal(t, expectedName, saveName)
				assert.True(t, ok)
			}
		})
	}

}
