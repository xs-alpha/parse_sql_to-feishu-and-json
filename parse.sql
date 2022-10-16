CREATE TABLE `parking_order`
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

CREATE TABLE `parking_peer_recorder`  (
  `recorder_id` bigint  NOT NULL AUTO_INCREMENT COMMENT '记录id',
  `order_id` varchar(64) NOT NULL COMMENT '订单号',
  `user_id` bigint NOT NULL COMMENT '乘客id',
  `passenger_job_number` varchar(20) NOT NULL DEFAULT '' COMMENT '乘客工号',
  `passenger_name` varchar(20) NOT NULL DEFAULT '' COMMENT '乘客姓名',
  `driver_job_number` varchar(20) NOT NULL DEFAULT '' COMMENT '司机工号',
  `driver_name` varchar(20) NOT NULL DEFAULT '' COMMENT '司机姓名',
  `set_out_place` varchar(50)  NOT NULL COMMENT '起始地',
  `dest_place` varchar(50)  NOT NULL COMMENT '目的地',
  `order_date` date NOT NULL COMMENT '出发日期',
  `order_time` time NOT NULL COMMENT '出发时间',
  `record_state` int NOT NULL DEFAULT 1 COMMENT '状态(
       1. 已请求：请求搭车，司机未点确认同行 
       2. 取消请求：乘客自己取消请求搭车
       3. 未成单： 司机确认同行其他乘客
       4. 已确认：司机确认同行
       5. 已取消：司机取消订单，过期自动取消订单
       6. 已结束：乘客到达目的地)',
  `create_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`recorder_id`)
) ENGINE = InnoDB CHARACTER SET = utf8 COMMENT = '同行记录表';