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

INSERT INTO Teams (Name, Strength, Points, GoalsFor, GoalsAgainst, MatchesPlayed)
VALUES 
    ('Arsenal', 85, 0, 0, 0, 0),
    ('Chelsea', 80, 0, 0, 0, 0),
    ('Liverpool', 90, 0, 0, 0, 0),
    ('Manchester United', 75, 0, 0, 0, 0);