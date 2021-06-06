create database if not exists shop ;

drop table if exists `users`;

create table if not exists `users` (
    `id` bigint not null auto_increment comment 'id',
    `name` varchar(128) not null default '' comment 'name',
    `passwd` varchar(64) not null default '79b11ef0ca0d0219eece85d3704d76a51a91a8e3a45bc17303cd8273667031a7' comment 'passwd sha256("passwd"+"1gsoOZ9bFyD7J2Q8K8hk6hrMk8H")',
    `sex` varchar(8) not null default '' comment 'sexual',
    `birthday` datetime default null comment 'birthday',
    `address` varchar(128) default null comment 'address',
    `email` varchar(128) default null comment 'email address',
    `version` varchar(64) not null default 'alpha 0.1' comment 'version',
    `is_delete` bool default false comment 'delete status',
    `gmt_create` datetime default current_timestamp comment 'create datetime',
    `gmt_update` datetime on update current_timestamp comment 'update datetime',
    primary key(`id`)
)engine=innodb default charset=utf8mb4