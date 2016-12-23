CREATE TABLE `noggo`.`admin` (
	`id` int NOT NULL AUTO_INCREMENT,
	`account` varchar(50) NOT NULL,
	`password` varchar(50) NOT NULL,
	`name` varchar(50) NOT NULL,
	`role` varchar(50) NOT NULL DEFAULT 'nobody',
	`changed` datetime NOT NULL DEFAULT now(),
	`created` datetime NOT NULL DEFAULT now(),
	PRIMARY KEY (`id`)
) COMMENT='';

INSERT INTO `noggo`.`admin`(`account`,`password`,`name`,`role`) VALUES ('admin','admin','admin','system');
