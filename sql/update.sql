/*
SQLyog Ultimate - MySQL GUI v8.2
MySQL - 5.5.27 : Database - hengzhu
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`hengzhu` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;

USE `hengzhu`;

/*Table structure for table `admin` */

DROP TABLE IF EXISTS `admin`;

CREATE TABLE `admin` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `login_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `real_name` varchar(32) NOT NULL DEFAULT '0' COMMENT '真实姓名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `role_ids` varchar(255) NOT NULL DEFAULT '0' COMMENT '角色id字符串，如：2,3,4',
  `phone` varchar(20) NOT NULL DEFAULT '0' COMMENT '手机号码',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，1-正常 0禁用',
  `create_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `update_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改者ID',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`login_name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

/*Data for the table `admin` */

insert  into `admin`(`id`,`login_name`,`real_name`,`password`,`role_ids`,`phone`,`email`,`salt`,`last_login`,`last_ip`,`status`,`create_id`,`update_id`,`create_time`,`update_time`) values (1,'admin','超管啊','4fd71bffda5ccb3d750931d764fd9979','1','13888888889','124@163.com','d9Fr',1512659112,'[',1,0,1,0,1512404296),(6,'test','测试普通管理员','a8dedcce538e590c79d7eb74358ea47a','2','','','3Kv3',1512659030,'[',1,1,0,1512658694,1512659041);

/*Table structure for table `auth` */

DROP TABLE IF EXISTS `auth`;

CREATE TABLE `auth` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `pid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上级ID，0为顶级',
  `auth_name` varchar(64) NOT NULL DEFAULT '0' COMMENT '权限名称',
  `auth_url` varchar(255) NOT NULL DEFAULT '0' COMMENT 'URL地址',
  `sort` int(1) unsigned NOT NULL DEFAULT '999' COMMENT '排序，越小越前',
  `icon` varchar(255) NOT NULL,
  `is_show` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否显示，0-隐藏，1-显示',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '操作者ID',
  `create_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `update_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改者ID',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态，1-正常，0-删除',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COMMENT='权限因子';

/*Data for the table `auth` */

insert  into `auth`(`id`,`pid`,`auth_name`,`auth_url`,`sort`,`icon`,`is_show`,`user_id`,`create_id`,`update_id`,`status`,`create_time`,`update_time`) values (1,0,'所有权限','/',1,'',0,1,1,1,1,1505620970,1505620970),(2,1,'权限管理','/',999,'fa-id-card',1,1,0,1,1,0,1505622360),(3,2,'管理员','/admin/list',1,'fa-user-o',1,1,1,1,1,1505621186,1505621186),(4,2,'角色管理','/role/list',2,'fa-user-circle-o',1,1,0,1,1,0,1505621852),(5,3,'新增','/admin/add',1,'',0,1,0,1,1,0,1505621685),(6,3,'修改','/admin/edit',2,'',0,1,0,1,1,0,1505621697),(7,3,'删除','/admin/ajaxdel',3,'',0,1,1,1,1,1505621756,1505621756),(8,4,'新增','/role/add',1,'',1,1,0,1,1,0,1505698716),(9,4,'修改','/role/edit',2,'',0,1,1,1,1,1505621912,1505621912),(10,4,'删除','/role/ajaxdel',3,'',0,1,1,1,1,1505621951,1505621951),(11,2,'权限因子','/auth/list',3,'fa-list',1,1,1,1,1,1505621986,1505621986),(12,11,'新增','/auth/add',1,'',0,1,1,1,1,1505622009,1505622009),(13,11,'修改','/auth/edit',2,'',0,1,1,1,1,1505622047,1505622047),(14,11,'删除','/auth/ajaxdel',3,'',0,1,1,1,1,1505622111,1505622111),(15,1,'个人中心','profile/edit',1001,'fa-user-circle-o',1,1,0,1,1,0,1506001114),(24,15,'资料修改','/user/edit',1,'fa-edit',1,1,0,1,1,0,1506057468),(39,1,'状态详情','cabinet',1,'fa-id-card',1,1,0,1,1,0,1511972856),(40,39,'状态列表','/cabinet/list',1,'th-list',1,1,0,1,1,0,1512396518),(41,40,'详情','/cabinet/detail',1,'',1,1,0,1,1,0,1512185124),(42,40,'详情列表','/cabinetDetail/table',2,'',1,1,0,1,1,0,1512196800),(43,40,'柜子启用','/cabinetdetail/changeuse',3,'th-list',1,1,0,1,1,0,1512234272),(44,39,'状态','/cabinetdetail/clear',4,'th-list',1,1,1,1,0,1512236453,1512236453),(45,40,'柜子清除','/cabinetdetail/clear',4,'',0,1,1,1,1,1512236513,1512236513),(46,40,'柜子记录','/cabinetdetail/record',5,'',1,1,0,1,1,0,1512309176),(47,39,'类型列表','/types/list',1,'th-list',1,1,1,1,1,1512396568,1512396568),(48,47,'设置默认','/types/default',1,'th-list',1,1,1,1,1,1512492807,1512492807),(49,47,'删除类型','/types/delete',2,'th-list',1,1,1,1,1,1512493301,1512493301),(50,47,'新增类型','/types/add',3,'th-list',1,1,1,1,1,1512573737,1512573737),(51,39,'设置','/setting/get',1,'th-list',1,0,0,0,1,0,1512661160);

/*Table structure for table `cabinet` */

DROP TABLE IF EXISTS `cabinet`;

CREATE TABLE `cabinet` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cabinet_ID` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '柜子id',
  `type_id` int(11) DEFAULT NULL COMMENT '柜子计费类型id，初始化时为默认类型',
  `address` text COLLATE utf8_unicode_ci COMMENT '柜子位置',
  `number` text COLLATE utf8_unicode_ci COMMENT '编号',
  `desc` text COLLATE utf8_unicode_ci COMMENT '备注',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_time` timestamp NULL DEFAULT NULL COMMENT '最后一次上报时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `cabinet` */

insert  into `cabinet`(`id`,`cabinet_ID`,`type_id`,`address`,`number`,`desc`,`create_time`,`last_time`,`update_time`) values (1,'012345',1,'给对方过后','asdf','asdf',NULL,'2017-11-25 15:29:00','2017-12-06 00:14:03'),(2,'123',2,NULL,NULL,NULL,'2017-11-24 22:17:13','2017-12-05 00:14:00',NULL);

/*Table structure for table `cabinet_detail` */

DROP TABLE IF EXISTS `cabinet_detail`;

CREATE TABLE `cabinet_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cabinet_id` int(11) DEFAULT NULL COMMENT '柜子的id',
  `door` int(11) DEFAULT NULL COMMENT '门号',
  `open_state` int(11) DEFAULT '1' COMMENT '开关状态，1:关，2:开',
  `using` int(11) DEFAULT '1' COMMENT '占用状态，1:空闲，2:占用',
  `userID` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '存物ID',
  `store_time` int(11) DEFAULT NULL COMMENT '存物时间',
  `use_state` int(11) DEFAULT '1' COMMENT '启用状态，1:启用，2:停用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `cabinet_detail` */

insert  into `cabinet_detail`(`id`,`cabinet_id`,`door`,`open_state`,`using`,`userID`,`store_time`,`use_state`) values (1,1,1,1,1,'',0,1),(2,1,2,1,1,'',0,1),(3,1,3,1,1,'',0,1);

/*Table structure for table `log` */

DROP TABLE IF EXISTS `log`;

CREATE TABLE `log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cabinet_detail_id` int(11) DEFAULT NULL COMMENT '柜子详情中的id',
  `action` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '操作',
  `user` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '存物ID/管理员',
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `log` */

insert  into `log`(`id`,`cabinet_detail_id`,`action`,`user`,`time`) values (1,1,'清除','超管','2017-12-03 01:42:42'),(2,1,'清除','超管','2017-12-03 21:47:06'),(3,1,'清除','超管啊','2017-12-05 23:46:02'),(4,1,'清除','超管啊','2017-12-06 00:06:33'),(5,1,'清除','超管啊','2017-12-06 00:07:01'),(6,3,'清除','超管啊','2017-12-06 00:11:50');

/*Table structure for table `role` */

DROP TABLE IF EXISTS `role`;

CREATE TABLE `role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_name` varchar(32) NOT NULL DEFAULT '0' COMMENT '角色名称',
  `detail` varchar(255) NOT NULL DEFAULT '0' COMMENT '备注',
  `create_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `update_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改这ID',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态1-正常，0-删除',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='角色表';

/*Data for the table `role` */

insert  into `role`(`id`,`role_name`,`detail`,`create_id`,`update_id`,`status`,`create_time`,`update_time`) values (1,'超级管理员','超级管理员，拥有所有权限',0,0,1,1512659848,1512659848),(2,'普通','拥有部分权限',0,1,1,1512657612,1512657612);

/*Table structure for table `role_auth` */

DROP TABLE IF EXISTS `role_auth`;

CREATE TABLE `role_auth` (
  `role_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
  `auth_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '权限ID',
  PRIMARY KEY (`role_id`,`auth_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限和角色关系表';

/*Data for the table `role_auth` */

insert  into `role_auth`(`role_id`,`auth_id`) values (1,0),(1,1),(1,2),(1,3),(1,4),(1,5),(1,6),(1,7),(1,8),(1,9),(1,10),(1,11),(1,12),(1,13),(1,14),(1,15),(1,24),(1,39),(1,40),(1,41),(1,42),(1,43),(1,45),(1,46),(1,47),(1,48),(1,49),(1,50),(1,51);

/*Table structure for table `setting` */

DROP TABLE IF EXISTS `setting`;

CREATE TABLE `setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `alipay` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '支付宝收款方',
  `weixin` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '微信收款方',
  `log_time` int(11) DEFAULT '30' COMMENT '历史记录保存时间',
  `customer` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '客服电话',
  `update_time` int(11) DEFAULT NULL COMMENT '创建/更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `setting` */

insert  into `setting`(`id`,`alipay`,`weixin`,`log_time`,`customer`,`update_time`) values (1,'','',15,'18412341234',1512664370);

/*Table structure for table `type` */

DROP TABLE IF EXISTS `type`;

CREATE TABLE `type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '类型名称',
  `default` int(11) DEFAULT NULL COMMENT '是否默认，1:默认，2:否',
  `charge_mode` int(11) DEFAULT '1' COMMENT '计费方式，1:计次，2:计时',
  `toll_time` int(11) DEFAULT '1' COMMENT '收费时间，1:存物时，2:取物时',
  `price` double DEFAULT NULL COMMENT '价格，若方式为计次，则价格为每次存取物价格，若方式为计时，则价格为unit时间内价格',
  `unit` int(11) DEFAULT NULL COMMENT '计时单位（分钟），当计费方式为计时时有',
  `create_time` int(11) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `type` */

insert  into `type`(`id`,`name`,`default`,`charge_mode`,`toll_time`,`price`,`unit`,`create_time`) values (2,'类型2',2,2,2,2,30,1512393048),(3,'类型3',2,1,1,0.5,NULL,1512393048),(4,'asf',1,1,1,12,0,1512578712),(5,'asdf',2,2,2,1,10,1512578823),(6,'123',2,2,2,234,234,1512579635),(7,'asdf',2,1,1,1,0,1512664263);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
