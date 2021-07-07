GRANT ALL PRIVILEGES ON DATABASE cryptos TO lachezar;

CREATE SCHEMA Cryptos;
GRANT ALL PRIVILEGES ON SCHEMA Cryptos TO lachezar;

CREATE TABLE Cryptos.Cryptocurrencies (
                                          Name varchar(20) NOT NULL,
                                          CryptoID varchar(10) NOT NULL
                                              CONSTRAINT PK_Cryptocurrencies_CryptoID PRIMARY KEY,
                                          Price numeric(10, 2) NOT NULL
                                              CONSTRAINT CK_Cryptocurrencies_Price_must_be_positive
                                                  CHECK (Price > 0)
);

CREATE TABLE Cryptos.Authors (
                                 CryptoID varchar(10) NULL
                                     CONSTRAINT FK_Authors_CryptoID
                                     REFERENCES Cryptos.Cryptocurrencies(CryptoID)
                                     ON DELETE CASCADE
                                     ON UPDATE CASCADE,
                                 CONSTRAINT PK_Authors PRIMARY KEY (CryptoID, firstname, lastname),
                                 Firstname varchar(20) NULL
                                     CONSTRAINT DK_Authors_Firstname_Unknown
                                     DEFAULT 'Unknown',
                                 Lastname varchar(20) NULL
                                     CONSTRAINT DK_Authors_Lastname_Unknown
                                     DEFAULT 'Unknown'
);

CREATE TABLE IF NOT EXISTS Cryptos.Cryptocurrencies_Audit (
    Name varchar(20) NOT NULL,
    CryptoID varchar(10) NOT NULL
    CONSTRAINT PK_Cryptocurrencies_Audit_CryptoID PRIMARY KEY,
    Price numeric(10, 2) NOT NULL
    CONSTRAINT CK_Cryptocurrencies_Audit_Price_must_be_positive
    CHECK (Price > 0),
    Doer varchar(20) NOT NULL,
    CryptoAdditionTime DATE
);


CREATE OR REPLACE FUNCTION Cryptocurrencies_Insert_Delete_Update_Fnc()
    RETURNS TRIGGER AS
    $$
    BEGIN
        IF (TG_OP = 'INSERT') THEN
            INSERT INTO Cryptos.Cryptocurrencies_Audit(Name, CryptoID, Price, Doer, CryptoAdditionTime)
            VALUES (NEW.Name, NEW.CryptoID, New.Price, current_user, current_date);

            RETURN NEW;
        ELSIF (TG_OP = 'DELETE') THEN
            DELETE FROM Cryptos.Cryptocurrencies_Audit WHERE CryptoID = OLD.CryptoID;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            UPDATE Cryptos.Cryptocurrencies_Audit SET Name = NEW.Name, CryptoID = NEW.CryptoID, Price = NEW.Price,
                                                      Doer = current_user, CryptoAdditionTime = current_date
            WHERE CryptoID = OLD.CryptoID;
            RETURN NEW;
        END IF;
    END;
    $$
    LANGUAGE plpgsql;

CREATE TRIGGER Cryptocurrencies_Insert_Delete_Update_Trigger
    AFTER INSERT OR DELETE OR UPDATE ON Cryptos.Cryptocurrencies
    FOR EACH ROW
    EXECUTE PROCEDURE Cryptocurrencies_Insert_Delete_Update_Fnc();

INSERT INTO Cryptos.Cryptocurrencies (
    Name,
    CryptoID,
    Price)
VALUES
('Bitcoin', 'BTC', '45000.94'),
('Ethereum', 'ETH', '2500.12'),
('DefaultCoin', 'DFC', '0.03');


INSERT INTO Cryptos.Authors (
    CryptoID,
    Firstname,
    Lastname)
VALUES
('BTC', 'Satoshi', 'Nakamoto'),
('ETH', 'Vitalik', 'Buterin'),
('ETH', 'Gavin', 'Wood');