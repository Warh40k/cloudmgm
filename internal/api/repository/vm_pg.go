package repository

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VmPostgres struct {
	db *sqlx.DB
}

func NewVmPostgres(db *sqlx.DB) *VmPostgres {
	return &VmPostgres{db: db}
}

func (r VmPostgres) ListVm(id uuid.UUID) ([]domain.VirtualMachine, error) {
	var vms []domain.VirtualMachine
	query := fmt.Sprintf(`SELECT v.id,v.title,v.description,v.created FROM %s v 
         INNER JOIN %s uv 
         ON v.id = uv.vm_id 
         WHERE uv.user_id = $1`, vmsTable, usersVmsTable)
	err := r.db.Select(&vms, query, id)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func (r VmPostgres) GetVm(id uuid.UUID) (domain.VirtualMachine, error) {
	//TODO implement me
	panic("implement me")
}

func (r VmPostgres) CreateVm(id uuid.UUID, machine domain.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func (r VmPostgres) DeleteVm(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r VmPostgres) ModifyVm(id uuid.UUID, machine domain.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}
