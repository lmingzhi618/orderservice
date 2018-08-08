CREATE TABLE orderserver.t_orders (
    `id`            INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `origin`        VARCHAR(30)  NOT NULL DEFAULT '',
    `destination`   VARCHAR(30)  NOT NULL DEFAULT '',
    `distance`      INT          NOT NULL DEFAULT 0,
    `status`        INT          NOT NULL DEFAULT 0 COMMENT '0: unassign, 1: order_have_been_taken'
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
