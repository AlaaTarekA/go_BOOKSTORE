
-- MySQL dump 10.13  Distrib 8.0.13, for macos10.14 (x86_64)
--
-- Host: localhost    Database: books_database
-- ------------------------------------------------------


--
-- Table structure for table `book_data`
--

CREATE TABLE `book_data` (
  `Title` varchar(255) NOT NULL,
  `Author` varchar(255) NOT NULL,
  `Publisher` varchar(255) NOT NULL,
  `PublishDate` datetime DEFAULT NULL,
  `Rating` int(64) DEFAULT NULL,
  `Status` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1

--
-- Dumping data for table `address_book`
--

LOCK TABLES `book_data` WRITE;
INSERT INTO `book_data` VALUES ("ketab","roba","solly","2019-10-07 00:00:00",2,"unchecked"),("abc", "mohannad", "pooBoo", "2019-11-01 00:00:00", 1, "checked");
UNLOCK TABLES;
