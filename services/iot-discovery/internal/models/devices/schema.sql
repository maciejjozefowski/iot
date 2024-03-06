create table devices
(
    id    serial primary key,
    name  text not null,
    brand text not null,
    model text not null,
    mac   text not null
);
