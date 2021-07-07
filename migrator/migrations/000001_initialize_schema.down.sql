DROP TRIGGER IF EXISTS Cryptocurrencies_Insert_Delete_Update_Trigger ON Cryptos.Cryptocurrencies;

DROP TABLE IF EXISTS Cryptos.Authors, Cryptos.Cryptocurrencies, Cryptos.Cryptocurrencies_Audit;

DROP SCHEMA IF EXISTS Cryptos;
