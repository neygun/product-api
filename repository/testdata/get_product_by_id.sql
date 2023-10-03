truncate table "product";
insert into "product" (id, name, price, created_at, updated_at) values (1, 'test', 1, now(), now());
select * from "product" where id=1;
