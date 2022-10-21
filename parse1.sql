CREATE TABLE `parking_orderu`
(
    `order_id`      varchar(64)  NOT NULL COMMENT '订单号',
    `user_id`       bigint       NOT NULL COMMENT '司机id',
    `job_number`    varchar(20)  NOT NULL DEFAULT '' COMMENT '工号',
    `name`          varchar(20)  NOT NULL DEFAULT '' COMMENT '姓名',
    `set_out_place` varchar(50)  NOT NULL COMMENT '起始地',
    `dest_place`    varchar(50)  NOT NULL COMMENT '目的地',
    `order_date`    date         NOT NULL COMMENT '出发日期',
    `order_time`    time         NOT NULL COMMENT '出发时间',
    `passenger_num` int          NOT NULL DEFAULT 1 COMMENT '可乘车人数',
    `order_comment` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
    `order_state`   int          NOT NULL DEFAULT 1 COMMENT '订单状态(1:未开始,2:进行中,3:已取消,4.已完成)',
    `create_at`     timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `update_at`     timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP (3) COMMENT '更新时间',
    PRIMARY KEY (`order_id`)
) ENGINE = InnoDB CHARACTER SET = UTF8MB4 COMMENT = '订单表';
