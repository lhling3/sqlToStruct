CREATE TABLE
    `commercial_use_record` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                `auth_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '授权id',
                                `sku_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '资源id',
                                `privilege_id` varchar(25) NOT NULL COMMENT '权益id',
                                `sku_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1:图片；2：图标；3：模板；4：字体',
                                `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
                                `nickname` varchar(255) DEFAULT '' COMMENT '名称',
                                `thumburl` varchar(255) DEFAULT '' COMMENT '资源缩略图',
                                `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
                                `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `idx_uk_user_id_sku_id` (`user_id`, `sku_id`, `sku_type`, `privilege_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 335 DEFAULT CHARSET = utf8mb4 COMMENT = '商用授权使用记录'