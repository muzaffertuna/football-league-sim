-- Teams tablosu
CREATE TABLE Teams (
    ID INT IDENTITY(1,1) PRIMARY KEY,
    Name NVARCHAR(100) NOT NULL,
    Strength INT NOT NULL,
    Points INT DEFAULT 0,
    GoalsFor INT DEFAULT 0,
    GoalsAgainst INT DEFAULT 0,
    MatchesPlayed INT DEFAULT 0
);

-- Matches tablosu
CREATE TABLE Matches (
    ID INT IDENTITY(1,1) PRIMARY KEY,
    HomeTeamID INT NOT NULL,
    AwayTeamID INT NOT NULL,
    HomeGoals INT DEFAULT 0,
    AwayGoals INT DEFAULT 0,
    Week INT NOT NULL,
    Played BIT DEFAULT 0,
    FOREIGN KEY (HomeTeamID) REFERENCES Teams(ID),
    FOREIGN KEY (AwayTeamID) REFERENCES Teams(ID)
);

-- Örnek takımlar
INSERT INTO Teams (Name, Strength, Points, GoalsFor, GoalsAgainst, MatchesPlayed)
VALUES 
    ('Arsenal', 85, 0, 0, 0, 0),
    ('Chelsea', 80, 0, 0, 0, 0),
    ('Liverpool', 90, 0, 0, 0, 0),
    ('Manchester United', 75, 0, 0, 0, 0);

-- Maç programı (4 takım, her biri diğerleriyle 2 kez oynar, toplam 12 maç, 6 hafta)
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (1, 2, 1, 0); -- Arsenal vs Chelsea
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (3, 4, 1, 0); -- Liverpool vs Man United
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (1, 3, 2, 0); -- Arsenal vs Liverpool
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (2, 4, 2, 0); -- Chelsea vs Man United
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (4, 1, 3, 0); -- Man United vs Arsenal
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (3, 2, 3, 0); -- Liverpool vs Chelsea
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (2, 1, 4, 0); -- Chelsea vs Arsenal
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (4, 3, 4, 0); -- Man United vs Liverpool
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (3, 1, 5, 0); -- Liverpool vs Arsenal
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (4, 2, 5, 0); -- Man United vs Chelsea
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (1, 4, 6, 0); -- Arsenal vs Man United
INSERT INTO Matches (HomeTeamID, AwayTeamID, Week, Played) VALUES (2, 3, 6, 0); -- Chelsea vs Liverpool