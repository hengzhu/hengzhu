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

/*Table structure for table `admin_user` */

DROP TABLE IF EXISTS `admin_user`;

CREATE TABLE `admin_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `nickname` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '2',
  `lastlogintime` datetime DEFAULT NULL,
  `createtime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `nickname` (`nickname`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `admin_user` */

insert  into `admin_user`(`id`,`username`,`password`,`nickname`,`email`,`remark`,`status`,`lastlogintime`,`createtime`) values (1,'admin','21232f297a57a5a743894a0e4a801fc3','ClownFish','osgochina@gmail.com','I\'m admin',2,NULL,'2017-11-20 07:37:32'),(2,'test','098f6bcd4621d373cade4e832627b4f6','test','23','',2,NULL,'2017-11-20 07:56:58');

/*Table structure for table `admin_user_roles` */

DROP TABLE IF EXISTS `admin_user_roles`;

CREATE TABLE `admin_user_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `admin_user_id` bigint(20) NOT NULL,
  `role_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `admin_user_roles` */

insert  into `admin_user_roles`(`id`,`admin_user_id`,`role_id`) values (1,2,2);

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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `cabinet` */

insert  into `cabinet`(`id`,`cabinet_ID`,`type_id`,`address`,`number`,`desc`,`create_time`,`last_time`) values (1,'012345',1,'爱的色放','asdf','备注',NULL,'2017-11-25 15:29:00'),(2,'123',2,NULL,NULL,NULL,'2017-11-24 22:17:13','2017-11-25 15:17:00');

/*Table structure for table `cabinet_detail` */

DROP TABLE IF EXISTS `cabinet_detail`;

CREATE TABLE `cabinet_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cabinet_id` int(11) DEFAULT NULL COMMENT '柜子的id',
  `door` int(11) DEFAULT NULL COMMENT '门号',
  `open_state` int(11) DEFAULT '1' COMMENT '开关状态，1:关，2:开',
  `using` int(11) DEFAULT '1' COMMENT '占用状态，1:空闲，2:占用',
  `userID` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '存物ID',
  `store_time` datetime DEFAULT NULL COMMENT '存物时间',
  `use_state` int(11) DEFAULT '1' COMMENT '启用状态，1:启用，2:停用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `cabinet_detail` */

insert  into `cabinet_detail`(`id`,`cabinet_id`,`door`,`open_state`,`using`,`userID`,`store_time`,`use_state`) values (1,1,1,1,2,'alipay-123','2017-11-25 10:23:23',1),(2,1,2,1,1,NULL,NULL,2);

/*Table structure for table `group` */

DROP TABLE IF EXISTS `group`;

CREATE TABLE `group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `title` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '2',
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `group` */

insert  into `group`(`id`,`name`,`title`,`status`,`sort`) values (1,'APP','System',2,1);

/*Table structure for table `log` */

DROP TABLE IF EXISTS `log`;

CREATE TABLE `log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cabinet_detail_id` int(11) DEFAULT NULL COMMENT '柜子详情中的id',
  `action` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '操作',
  `user` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '存物ID/管理员',
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `log` */

/*Table structure for table `logs` */

DROP TABLE IF EXISTS `logs`;

CREATE TABLE `logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `action` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status_code` int(11) DEFAULT NULL,
  `input` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `logs` */

/*Table structure for table `node` */

DROP TABLE IF EXISTS `node`;

CREATE TABLE `node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `level` int(11) NOT NULL DEFAULT '1',
  `pid` bigint(20) NOT NULL DEFAULT '0',
  `remark` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '2',
  `group_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `node` */

insert  into `node`(`id`,`title`,`name`,`level`,`pid`,`remark`,`status`,`group_id`) values (1,'RBAC','rbac',1,0,'',2,1),(2,'Node','node/index',2,1,'',2,1),(3,'node list','index',3,2,'',2,1),(4,'add or edit','AddAndEdit',3,2,'',2,1),(5,'del node','DelNode',3,2,'',2,1),(6,'User','user/index',2,1,'',2,1),(7,'user list','Index',3,6,'',2,1),(8,'add user','AddUser',3,6,'',2,1),(9,'update user','UpdateUser',3,6,'',2,1),(10,'del user','DelUser',3,6,'',2,1),(11,'Group','group/index',2,1,'',2,1),(12,'group list','index',3,11,'',2,1),(13,'add group','AddGroup',3,11,'',2,1),(14,'update group','UpdateGroup',3,11,'',2,1),(15,'del group','DelGroup',3,11,'',2,1),(16,'Role','role/index',2,1,'',2,1),(17,'role list','index',3,16,'',2,1),(18,'add or edit','AddAndEdit',3,16,'',2,1),(19,'del role','DelRole',3,16,'',2,1),(20,'get roles','Getlist',3,16,'',2,1),(21,'show access','AccessToNode',3,16,'',2,1),(22,'add accsee','AddAccess',3,16,'',2,1),(23,'show role to userlist','RoleToUserList',3,16,'',2,1),(24,'add role to user','AddRoleToUser',3,16,'',2,1),(25,'状态','state/index',1,0,'',2,1);

/*Table structure for table `node_roles` */

DROP TABLE IF EXISTS `node_roles`;

CREATE TABLE `node_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `node_id` bigint(20) NOT NULL,
  `role_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `node_roles` */

insert  into `node_roles`(`id`,`node_id`,`role_id`) values (13,1,2),(14,16,2),(15,17,2),(16,18,2),(17,19,2),(18,20,2),(19,21,2),(20,22,2),(21,23,2),(22,24,2),(23,25,2),(24,26,2);

/*Table structure for table `role` */

DROP TABLE IF EXISTS `role`;

CREATE TABLE `role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `role` */

insert  into `role`(`id`,`title`,`name`,`remark`,`status`) values (1,'Admin role','Admin','I\'m a admin role',2),(2,'','Manager','',2);

/*Table structure for table `setting` */

DROP TABLE IF EXISTS `setting`;

CREATE TABLE `setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `alipay` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '支付宝收款方',
  `weixin` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '微信收款方',
  `log_time` int(11) DEFAULT '30' COMMENT '历史记录保存时间',
  `customer` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '客服电话',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `setting` */

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
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `type` */

insert  into `type`(`id`,`name`,`default`,`charge_mode`,`toll_time`,`price`,`unit`,`create_time`) values (1,'类型1',1,1,1,1,NULL,'2017-11-25 15:24:09'),(2,'类型2',2,2,2,2,30,'2017-11-25 15:28:53');

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `mail` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `mobile` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `nick_name` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0',
  `user_type` tinyint(4) NOT NULL DEFAULT '0',
  `register_wlan` varchar(17) COLLATE utf8_unicode_ci DEFAULT NULL,
  `game_id` int(11) NOT NULL DEFAULT '0',
  `account_mark` tinyint(4) NOT NULL DEFAULT '0',
  `initialpwd` varchar(12) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `solidify` tinyint(4) NOT NULL DEFAULT '0',
  `chinese_name` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `id_card_no` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `user` */

/*Table structure for table `user_bindmobile` */

DROP TABLE IF EXISTS `user_bindmobile`;

CREATE TABLE `user_bindmobile` (
  `bind_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL DEFAULT '0',
  `mobile` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `seccode` varchar(6) COLLATE utf8_unicode_ci DEFAULT NULL,
  `create_time` int(11) NOT NULL DEFAULT '0',
  `send_time` int(10) unsigned NOT NULL DEFAULT '0',
  `seccode_time` int(10) unsigned NOT NULL DEFAULT '0',
  `expire_time` int(11) NOT NULL DEFAULT '0',
  `mark` tinyint(4) DEFAULT NULL,
  `send_times` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`bind_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `user_bindmobile` */

/*Table structure for table `user_info` */

DROP TABLE IF EXISTS `user_info`;

CREATE TABLE `user_info` (
  `user_id` int(11) NOT NULL,
  `nick_name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `sex` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `birth` date NOT NULL,
  `headpic` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  `real_name` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `identity_card` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_tel` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_mobile` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_postcode` varchar(6) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_job` tinyint(4) NOT NULL DEFAULT '0',
  `user_salary` tinyint(4) NOT NULL DEFAULT '0',
  `user_education` tinyint(4) NOT NULL DEFAULT '0',
  `user_marital` tinyint(4) NOT NULL DEFAULT '0',
  `area` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `exp` int(10) unsigned NOT NULL DEFAULT '0',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `nosign` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*Data for the table `user_info` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
