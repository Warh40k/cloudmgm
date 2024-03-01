package postgres

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VolumePostgres struct {
	db *sqlx.DB
}

func (r VolumePostgres) CheckOwnership(userId, volumeId uuid.UUID) error {
	var vm domain.Volume
	query := fmt.Sprintf(`SELECT vm.id FROM %s vm 
								INNER JOIN %s uv ON vm.id = uv.volume_id 
								WHERE uv.user_id = $1 and uv.volume_id = $2`, volumesTable, usersVolumesTable)

	return r.db.Get(&vm, query, userId, volumeId)
}

func NewVolumePostgres(db *sqlx.DB) *VolumePostgres {
	return &VolumePostgres{db: db}
}

func (r VolumePostgres) ListVolume(userId uuid.UUID) ([]domain.Volume, error) {
	var volumes []domain.Volume
	query := fmt.Sprintf(`SELECT v.id,v.label,v.description, v.size, v.created FROM %s v 
         INNER JOIN %s uv 
         ON v.id = uv.volume_id 
         WHERE uv.user_id = $1`, volumesTable, usersVolumesTable)
	err := r.db.Select(&volumes, query, userId)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

func (r VolumePostgres) GetVolume(vmId uuid.UUID) (domain.Volume, error) {
	var vm domain.Volume
	query := fmt.Sprintf(`SELECT * FROM %s vms where vms.id = $1`, volumesTable)
	if err := r.db.Get(&vm, query, vmId); err != nil {
		return vm, ErrNoRows
	}

	return vm, nil
}

func (r VolumePostgres) CreateVolume(userId uuid.UUID, machine domain.Volume) (uuid.UUID, error) {
	var id uuid.UUID

	volumeId := uuid.New()
	tx, err := r.db.Beginx()

	if err != nil {
		return uuid.Nil, ErrInternal
	}
	vmQuery := fmt.Sprintf(`INSERT INTO %s(id,label,description) 
								VALUES($1,$2,$3) RETURNING id`, volumesTable)
	row := tx.QueryRowx(vmQuery, volumeId, machine.Label, machine.Description)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, ErrInternal
	}

	userVmId := uuid.New()
	userVmQuery := fmt.Sprintf(`INSERT INTO %s(id,user_id, volume_id) 
								VALUES($1,$2,$3)`, usersVolumesTable)
	_, err = tx.Exec(userVmQuery, userVmId, userId, volumeId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, ErrInternal
	}

	return id, tx.Commit()
}

func (r VolumePostgres) DeleteVolume(vmId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s v WHERE v.id = $1`, volumesTable)
	_, err := r.db.Exec(query, vmId)
	if err != nil {
		return err
	}

	return nil
}

func (r VolumePostgres) UpdateVolume(machine domain.Volume) error {
	query := fmt.Sprintf(`UPDATE %s 
								SET label = $1, description = $2 
								WHERE id = $3`, volumesTable)
	res, err := r.db.Exec(query, machine.Label, machine.Description, machine.Id)

	if err != nil {
		return ErrInternal
	}

	count, err := res.RowsAffected()
	if err != nil {
		return ErrInternal
	}

	if count == 0 {
		return ErrNoRows
	}

	return nil
}
