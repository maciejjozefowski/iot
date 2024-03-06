-- name: GetDeviceByMac :one
select * from devices where id = $1 limit 1;

-- name: GetDeviceByID :one
select * from devices where id = $1 limit 1;

-- name: GetDevicesListPaged :many
select * from devices order by id limit $1 offset $2;

-- name: GetDeviceListByNames :many
select * from devices where name = $1;

-- name: GetDeviceListByBrand :many
select * from devices where brand = $1;

-- name: GetDeviceListByModel :many
select * from devices where model = $1;

-- name: GetDeviceByParams :many
select * from devices where name = $1 or brand = $2 or model = $3 or mac = $4;

-- name: CreateDevice :one
insert into devices (name, brand, model, mac) values ($1, $2, $3, $4) returning *;

-- name: UpdateDevice :one
update devices set name = $1, brand = $2, model = $3, mac = $4 where id = $5 returning *;

-- name: DeleteDevice :one
delete from devices where id = $1 returning *;

