truncate table "product";
insert into "product" (id, name, price, created_at, updated_at) values (1, 'test1', 1, now(), now());
insert into "product" (id, name, price, created_at, updated_at) values (2, 'test2', 2, now(), now());
select * from "product";
