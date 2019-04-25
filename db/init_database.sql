create database `orderserver` default character set utf8 collate utf8_general_ci;
use orderserver;
drop table if exists `t_orders`;

CREATE TABLE `t_orders` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`origin` varchar(30) DEFAULT NULL,
`destination` varchar(30) DEFAULT NULL,
`distance` int(11) DEFAULT 0,
`status` int(11) DEFAULT 0,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
