
-- cannot init new db and create table in that db that is why I am using postgres db;
create database syook_timeseries;
create table TimeSeries(
    id text primary key,
    listenerId text,
    timeMinute timestamp default date_trunc('minute', now()),
    createdAt timestamp default now(),
    data text
);
create index ts_ld_tm on TimeSeries using btree (listenerId, timeminute);
create index ts_tm on TimeSeries using btree (timeminute);
