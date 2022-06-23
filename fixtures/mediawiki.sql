DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `comment_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `comment_hash` int NOT NULL,
  `comment_text` blob NOT NULL,
  `comment_data` blob,
  PRIMARY KEY (`comment_id`),
  KEY `comment_hash` (`comment_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;

INSERT INTO `comment` VALUES (1,0,'',NULL),(2,314159265,_binary 'A comment',NULL);

DROP TABLE IF EXISTS `page`;
CREATE TABLE `page` (
  `page_id` int unsigned NOT NULL AUTO_INCREMENT,
  `page_namespace` int NOT NULL,
  `page_title` varbinary(255) NOT NULL,
  `page_restrictions` tinyblob,
  `page_is_redirect` tinyint unsigned NOT NULL DEFAULT '0',
  `page_is_new` tinyint unsigned NOT NULL DEFAULT '0',
  `page_random` double unsigned NOT NULL,
  `page_touched` binary(14) NOT NULL,
  `page_latest` int unsigned NOT NULL,
  `page_len` int unsigned NOT NULL,
  `page_content_model` varbinary(32) DEFAULT NULL,
  `page_links_updated` varbinary(14) DEFAULT NULL,
  `page_lang` varbinary(35) DEFAULT NULL,
  PRIMARY KEY (`page_id`),
  UNIQUE KEY `page_name_title` (`page_namespace`,`page_title`),
  KEY `page_random` (`page_random`),
  KEY `page_len` (`page_len`),
  KEY `page_redirect_namespace_len` (`page_is_redirect`,`page_namespace`,`page_len`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;

INSERT INTO `page` VALUES
(2,0,_binary 'Page1','',0,1,0.87780065007 ,_binary '20220622120000',1,5,NULL,NULL,NULL),
(3,0,_binary 'Page2','',0,1,0.604584110411,_binary '20220622120000',2,5,NULL,NULL,NULL),
(4,0,_binary 'Page3','',0,0,0.520852071594,_binary '20220622130000',9,5,NULL,NULL,NULL),
(5,0,_binary 'Page4','',0,1,0.648016511812,_binary '20220622120000',4,5,NULL,NULL,NULL),
(6,0,_binary 'Page5','',0,1,0.279926780118,_binary '20220622120000',5,5,NULL,NULL,NULL),
(7,0,_binary 'Page6','',0,1,0.045127100005,_binary '20220622120000',6,5,NULL,NULL,NULL),
(8,0,_binary 'Page7','',0,1,0.716324798279,_binary '20220622120000',7,5,NULL,NULL,NULL);

DROP TABLE IF EXISTS `revision`;
CREATE TABLE `revision` (
  `rev_id` int unsigned NOT NULL AUTO_INCREMENT,
  `rev_page` int unsigned NOT NULL,
  `rev_comment_id` bigint unsigned NOT NULL DEFAULT '0',
  `rev_actor` bigint unsigned NOT NULL DEFAULT '0',
  `rev_timestamp` binary(14) NOT NULL,
  `rev_minor_edit` tinyint unsigned NOT NULL DEFAULT '0',
  `rev_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `rev_len` int unsigned DEFAULT NULL,
  `rev_parent_id` int unsigned DEFAULT NULL,
  `rev_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`rev_id`),
  KEY `rev_timestamp` (`rev_timestamp`),
  KEY `rev_page_id` (`rev_page`,`rev_id`),
  KEY `rev_actor_timestamp` (`rev_actor`,`rev_timestamp`,`rev_id`),
  KEY `rev_page_actor_timestamp` (`rev_page`,`rev_actor`,`rev_timestamp`),
  KEY `rev_page_timestamp` (`rev_page`,`rev_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary MAX_ROWS=10000000 AVG_ROW_LENGTH=1024;

INSERT INTO `revision` VALUES
(1,2,0,0,_binary '20220622120000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(2,3,0,0,_binary '20220622120000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(3,4,0,0,_binary '20220622120000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(4,5,0,0,_binary '20220622120000',1,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(5,6,0,0,_binary '20220622120000',1,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(6,7,0,0,_binary '20220622120000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(7,8,0,0,_binary '20220622120000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(8,4,0,0,_binary '20220622123000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234'),
(9,4,0,0,_binary '20220622130000',0,0,5,0,_binary 'abcdefghijklmnopqrstuvwxyz01234');

DROP TABLE IF EXISTS `revision_actor_temp`;
CREATE TABLE `revision_actor_temp` (
  `revactor_rev` int unsigned NOT NULL,
  `revactor_actor` bigint unsigned NOT NULL,
  `revactor_timestamp` binary(14) NOT NULL,
  `revactor_page` int unsigned NOT NULL,
  PRIMARY KEY (`revactor_rev`,`revactor_actor`),
  UNIQUE KEY `revactor_rev` (`revactor_rev`),
  KEY `actor_timestamp` (`revactor_actor`,`revactor_timestamp`),
  KEY `page_actor_timestamp` (`revactor_page`,`revactor_actor`,`revactor_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;

INSERT INTO `revision_actor_temp` VALUES
(1,1,_binary '20220622120000',2),
(2,2,_binary '20220622120000',3),
(3,3,_binary '20220622120000',4),
(4,4,_binary '20220622120000',5),
(5,5,_binary '20220622120000',6),
(6,6,_binary '20220622120000',7),
(7,7,_binary '20220622120000',8),
(8,8,_binary '20220622123000',4),
(9,9,_binary '20220622130000',4);

DROP TABLE IF EXISTS `revision_comment_temp`;
CREATE TABLE `revision_comment_temp` (
  `revcomment_rev` int unsigned NOT NULL,
  `revcomment_comment_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`revcomment_rev`,`revcomment_comment_id`),
  UNIQUE KEY `revcomment_rev` (`revcomment_rev`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;

INSERT INTO `revision_comment_temp` VALUES
(1,1),(2,1),(3,1),(4,2),(5,2),(6,1),(7,1),(8,2),(9,1);
