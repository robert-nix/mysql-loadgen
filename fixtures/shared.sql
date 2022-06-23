DROP TABLE IF EXISTS `actor`;
CREATE TABLE `actor` (
  `actor_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `actor_user` int unsigned DEFAULT NULL,
  `actor_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL,
  PRIMARY KEY (`actor_id`),
  UNIQUE KEY `actor_name` (`actor_name`),
  UNIQUE KEY `actor_user` (`actor_user`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
INSERT INTO `actor` (`actor_id`, `actor_user`, `actor_name`) VALUES
(1, 1, 'User1'),
(2, 2, 'User2'),
(3, 3, 'User3'),
(4, 4, 'User4'),
(5, 5, 'User5'),
(6, 6, 'User6'),
(7, 7, 'User7'),
(8, NULL, '127.0.0.2'),
(9, NULL, '127.0.0.3');

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL DEFAULT '',
  `user_real_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL DEFAULT '',
  `user_email` tinytext NOT NULL,
  `user_touched` char(14) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL DEFAULT '',
  `user_token` char(32) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT '',
  `user_email_authenticated` char(14) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
  `user_email_token` char(32) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
  `user_email_token_expires` char(14) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
  `user_registration` varchar(16) DEFAULT NULL,
  `user_editcount` int DEFAULT NULL,
  `user_birthdate` date DEFAULT NULL,
  `user_options` blob NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`),
  KEY `user_email_token` (`user_email_token`),
  KEY `user_email` (`user_email`(40)),
  KEY `user_registration` (`user_registration`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
INSERT INTO `user` (`user_id`, `user_name`, `user_real_name`, `user_email`, `user_touched`, `user_token`, `user_email_authenticated`, `user_email_token`, `user_email_token_expires`, `user_registration`, `user_editcount`, `user_birthdate`, `user_options`) VALUES
( 1, 'User1', '', 'user1@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', ''),
( 2, 'User2', '', 'user2@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', ''),
( 3, 'User3', '', 'user3@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', ''),
( 4, 'User4', '', 'user4@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', ''),
( 5, 'User5', '', 'user5@example.com' , '20120622133100' , '' , '20120622133100' , NULL , NULL , '20120622133100', 1 , '1970-01-01', ''),
( 6, 'User6', '', 'user6@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', ''),
( 7, 'User7', '', 'user7@example.com' , '20220622133100' , '' , '20220622133100' , NULL , NULL , '20220622133100', 1 , '1970-01-01', '');
