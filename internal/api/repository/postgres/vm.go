package postgres

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VmPostgres struct {
	db *sqlx.DB
}

func (r VmPostgres) CheckOwnership(userId, vmId uuid.UUID) error {
	var vm domain.VirtualMachine
	query := fmt.Sprintf(`SELECT vm.id FROM %s vm 
								INNER JOIN %s uv ON vm.id = uv.vm_id 
								WHERE uv.user_id = $1 and uv.vm_id = $2`, vmsTable, usersVmsTable)

	return r.db.Get(&vm, query, userId, vmId)
}

func NewVmPostgres(db *sqlx.DB) *VmPostgres {
	return &VmPostgres{db: db}
}

func (r VmPostgres) ListVm(userId uuid.UUID) ([]domain.VirtualMachine, error) {
	var vms []domain.VirtualMachine
	query := fmt.Sprintf(`SELECT v.id,v.title,v.description,v.created FROM %s v 
         INNER JOIN %s uv 
         ON v.id = uv.vm_id 
         WHERE uv.user_id = $1`, vmsTable, usersVmsTable)
	err := r.db.Select(&vms, query, userId)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func (r VmPostgres) GetVm(vmId uuid.UUID) (domain.VirtualMachine, error) {
	//TODO implement me
	panic("implement me")
}

func (r VmPostgres) CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error {
	vmId, err := uuid.NewUUID()
	if err != nil {
		return ErrInternal
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return ErrInternal
	}
	vmQuery := fmt.Sprintf(`INSERT INTO %s(id,title,description,status) 
								VALUES($1,$2,$3,0)`, vmsTable)
	_, err = tx.Exec(vmQuery, vmId, machine.Label, machine.Description)
	if err != nil {
		tx.Rollback()
		return ErrInternal
	}
	userVmId, err := uuid.NewUUID()
	if err != nil {
		tx.Rollback()
		return ErrInternal
	}
	userVmQuery := fmt.Sprintf(`INSERT INTO %s(id,user_id, vm_id) 
								VALUES($1,$2,$3)`, usersVmsTable)
	_, err = tx.Exec(userVmQuery, userVmId, userId, vmId)
	if err != nil {
		tx.Rollback()
		return ErrInternal
	}

	return tx.Commit()
}

func (r VmPostgres) DeleteVm(vmId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s v WHERE v.id = $1`, vmsTable)
	_, err := r.db.Exec(query, vmId)
	if err != nil {
		return err
	}

	return nil
}

func (r VmPostgres) ModifyVm(vmId uuid.UUID, machine domain.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}
