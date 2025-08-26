-- +goose Up

CREATE TABLE `OrderList` (
    `Id` varchar(36) NOT NULL,
    `Created` timestamp NOT NULL,
    `Deadline` timestamp,
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
    `DrinkId` varchar(36),
    `EditKey` text,
    PRIMARY KEY (`Id`),
    FOREIGN KEY (`OrderListId`) REFERENCES `OrderList`(`Id`) ON DELETE CASCADE,
    FOREIGN KEY (`DrinkId`) REFERENCES `Drink`(`Id`) ON DELETE CASCADE
);

CREATE TABLE `OrderItems` (
    `OrderId` varchar(36) NOT NULL,
    `StoreItemId` text NOT NULL,
    PRIMARY KEY (`OrderId`, `StoreItemId`),
    FOREIGN KEY (`OrderId`) REFERENCES `Order`(`Id`) ON DELETE CASCADE
);

CREATE TABLE `StoreItemVariant` (
    `OrderId` varchar(36) NOT NULL,
    `StoreItemId` text NOT NULL,
    `Variant` text NOT NULL,
    PRIMARY KEY (`OrderId`, `StoreItemId`, `Variant`),
    FOREIGN KEY (`OrderId`) REFERENCES `Order`(`Id`) ON DELETE CASCADE
);

CREATE TABLE `StoreItemDip` (
    `OrderId` varchar(36) NOT NULL,
    `StoreItemId` text NOT NULL,
    `Dip` text NOT NULL,
    PRIMARY KEY (`OrderId`, `StoreItemId`, `Dip`),
    FOREIGN KEY (`OrderId`) REFERENCES `Order`(`Id`) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE StoreItemDip;
DROP TABLE StoreItemVariant;
DROP TABLE OrderItems;
DROP TABLE `Order`;
DROP TABLE Drink;
DROP TABLE OrderList;
