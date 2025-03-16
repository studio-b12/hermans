-- +goose Up

CREATE TABLE `OrderList` (
    `Id` varchar(36) NOT NULL,
    `Created` timestamp NOT NULL,

    PRIMARY KEY (`Id`)
);

CREATE TABLE `Drink` (
     `Id` varchar(36) NOT NULL,
     `Name` text NOT NULL,
     `Size` int NOT NULL,

     PRIMARY KEY (`Id`)
 );

 CREATE TABLE `Order` (
    `Id` varchar(36) NOT NULL,
    `Created` timestamp NOT NULL,
    `Creator` text NOT NULL,
    `OrderListId` varchar(36) NOT NULL,
    `StoreItemId` text NOT NULL,
    `DrinkId` varchar(36),

    PRIMARY KEY (`Id`),
    FOREIGN KEY (`OrderListId`) REFERENCES `OrderList`(`Id`) ON DELETE CASCADE,
    FOREIGN KEY (`DrinkId`) REFERENCES `Drink`(`Id`) ON DELETE CASCADE
);

CREATE TABLE `StoreItemVariant` (
    `OrderId` varchar(36) NOT NULL,
    `Variant` text NOT NULL,

    PRIMARY KEY (`OrderId`, `Variant`),
    FOREIGN KEY (`OrderId`) REFERENCES `Order`(`Id`) ON DELETE CASCADE
);

CREATE TABLE `StoreItemDip` (
    `OrderId` varchar(36) NOT NULL,
    `Dip` text NOT NULL,

    PRIMARY KEY (`OrderId`, `Dip`),
    FOREIGN KEY (`OrderId`) REFERENCES `Order`(`Id`) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE StoreItemDip;
DROP TABLE StoreItemVariant;
DROP TABLE Order;
DROP TABLE Drink;
DROP TABLE OrderList;